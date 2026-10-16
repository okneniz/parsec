package lang

import (
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"
)

func TestLexerIndentation(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "a nested function",
			src:  "def f(x):\n    if x:\n        return 1\n    return 0\n",
			want: []string{
				"def", "f", "(", "x", ")", ":", "newline",
				"indent",
				"if", "x", ":", "newline",
				"indent",
				"return", "1", "newline",
				"dedent",
				"return", "0", "newline",
			},
		},
		{
			name: "blank lines carry no indentation",
			src:  "x = 1\n\ny = 2\n",
			want: []string{"x", "=", "1", "newline", "newline", "y", "=", "2", "newline"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// a fresh state per case: the indent stack belongs to one
			// program
			stream, err := lex(pythonDefinition(), newIndentState(), test.src)
			require.NoError(t, err)
			require.Equal(t, test.want, lexemeKinds(t, stream))
		})
	}
}

func TestLexerIndentationErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"unindent to no outer level", "if x:\n    y\n  z\n", "unindent does not match"},
		{"unindent between two blocks", "x:\n    y:\n        z\n  w\n", "unindent does not match"},
		{"unterminated string", "x = \"abc\n", "unterminated string"},
		{"stray character", "x = $\n", "unrecognized character"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := lex(pythonDefinition(), newIndentState(), test.src)
			require.ErrorContains(t, err, test.want)
		})
	}
}

// indentState is the lexer state of an indentation language: the
// stack of the open block widths and the dedents queued ahead of the
// tokens that follow them.
type indentState struct {
	stack   []int
	pending []Token
}

func (s *indentState) Clone() *indentState {
	clone := indentState{
		stack:   make([]int, 0, len(s.stack)),
		pending: make([]Token, 0, len(s.pending)),
	}

	clone.stack = append(clone.stack, s.stack...)
	clone.pending = append(clone.pending, s.pending...)

	return &clone
}

func (s *indentState) Recover(clone *indentState) {
	s.stack = clone.stack
	s.pending = clone.pending
}

func newIndentState() *indentState {
	return &indentState{stack: []int{0}}
}

// pythonDefinition describes a Python-flavored language: the newline
// is a token and the indent stack lives in the lexer state.
func pythonDefinition() Definition[*indentState] {
	return Definition[*indentState]{
		// the newline stays out of the trivia: it is a token
		Space: func(r rune) bool {
			return r == ' ' || r == '\t'
		},
		Keywords:    []string{"def", "if", "else", "return"},
		Operators:   []string{"=", "+", "-"},
		Punctuation: "(),:",
		Integers:    true,
		Strings:     true,
		Custom:      []func(state *indentState, buf parsec.Buffer[rune, strings.Position]) (Token, error){indentTokens},
	}
}

func lexemeKinds(t *testing.T, stream []Token) []string {
	t.Helper()

	out := make([]string, 0, len(stream))

	for _, tok := range stream {
		switch string(tok.Kind) {
		case "newline", "indent", "dedent":
			out = append(out, string(tok.Kind))
		default:
			out = append(out, string(tok.Lexeme))
		}
	}

	return out
}

// indentTokens reads the newline and indentation tokens of a
// Python-flavored language: a newline token, then the leading width
// of the next line against the stack — wider opens an indent,
// narrower closes blocks, a blank line only ends its own line. The
// dedents of one line queue ahead of the tokens that follow, which
// the customs-are-tried-every-token contract allows. The end of
// input ends the last line without closing dedents: the next call
// meets the end of input.
func indentTokens(state *indentState, buf parsec.Buffer[rune, strings.Position]) (Token, error) {
	if n := len(state.pending); n > 0 {
		tok := state.pending[0]
		state.pending = state.pending[1:]

		return tok, nil
	}

	r, err := buf.Read(false)
	if err != nil || r != '\n' {
		return Token{}, ErrNoMatch
	}

	pos := buf.Position()

	if _, err := buf.Read(true); err != nil {
		return Token{}, err
	}

	// the width of the next line; the spaces are trivia anyway,
	// consuming them here only takes what the stepper would
	width := 0

	for {
		at := buf.Position()

		r, err = buf.Read(true)
		if err != nil {
			break
		}

		if r != ' ' && r != '\t' {
			if err := buf.Seek(at); err != nil {
				return Token{}, err
			}

			break
		}

		width++
	}

	newline := Token{Kind: "newline", Lexeme: "\n"}

	// a blank line, or the end of input, ends its own line only
	peek, perr := buf.Read(false)
	if perr != nil || peek == '\n' {
		return newline, nil
	}

	top := state.stack[len(state.stack)-1]

	switch {
	case width > top:
		state.stack = append(state.stack, width)
		state.pending = append(state.pending, Token{Kind: "indent"})
	case width < top:
		for len(state.stack) > 1 && state.stack[len(state.stack)-1] > width {
			state.stack = state.stack[:len(state.stack)-1]
			state.pending = append(state.pending, Token{Kind: "dedent"})
		}

		if state.stack[len(state.stack)-1] != width {
			return Token{}, parsec.NewParseError(pos, "unindent does not match any outer level")
		}
	}

	return newline, nil
}
