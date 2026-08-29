package lang

import (
	"os"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"
)

func TestLexerDefinitions(t *testing.T) {
	tests := []struct {
		name string
		def  Definition
		src  string
		want []Token
	}{
		{
			name: "bare identifiers and integers",
			def:  Definition{Integers: true},
			src:  "x1 42 y",
			want: []Token{
				{Kind: KindIdent, Lexeme: "x1"},
				{Kind: KindInt, Lexeme: "42"},
				{Kind: KindIdent, Lexeme: "y"},
			},
		},
		{
			name: "a config file: line comments, keywords, punctuation",
			def: Definition{
				LineComment: "#",
				Keywords:    []string{"set", "on", "off"},
				Punctuation: "=:",
				Integers:    true,
			},
			src: "# note\nset debug = 1\nset mode : off",
			want: []Token{
				{Kind: KindKeyword, Lexeme: "set"},
				{Kind: KindIdent, Lexeme: "debug"},
				{Kind: KindSymbol, Lexeme: "="},
				{Kind: KindInt, Lexeme: "1"},
				{Kind: KindKeyword, Lexeme: "set"},
				{Kind: KindIdent, Lexeme: "mode"},
				{Kind: KindSymbol, Lexeme: ":"},
				{Kind: KindKeyword, Lexeme: "off"},
			},
		},
		{
			name: "shared operator prefixes go to the longest",
			def: Definition{
				Operators:   []string{"<<=", "<<", "<"},
				Punctuation: "()",
			},
			src: "a <<= b << c < (d)",
			want: []Token{
				{Kind: KindIdent, Lexeme: "a"},
				{Kind: KindOperator, Lexeme: "<<="},
				{Kind: KindIdent, Lexeme: "b"},
				{Kind: KindOperator, Lexeme: "<<"},
				{Kind: KindIdent, Lexeme: "c"},
				{Kind: KindOperator, Lexeme: "<"},
				{Kind: KindSymbol, Lexeme: "("},
				{Kind: KindIdent, Lexeme: "d"},
				{Kind: KindSymbol, Lexeme: ")"},
			},
		},
		{
			name: "nested comments, strings with escapes",
			def: Definition{
				BlockComment:   [2]string{"(*", "*)"},
				NestedComments: true,
				Operators:      []string{"==", "=", "<="},
				Punctuation:    "()",
				Integers:       true,
				Strings:        true,
			},
			src: "x <= y (* c (* d *) e *) == (z = \"a\\nb\")",
			want: []Token{
				{Kind: KindIdent, Lexeme: "x"},
				{Kind: KindOperator, Lexeme: "<="},
				{Kind: KindIdent, Lexeme: "y"},
				{Kind: KindOperator, Lexeme: "=="},
				{Kind: KindSymbol, Lexeme: "("},
				{Kind: KindIdent, Lexeme: "z"},
				{Kind: KindOperator, Lexeme: "="},
				{Kind: KindString, Lexeme: "a\nb"},
				{Kind: KindSymbol, Lexeme: ")"},
			},
		},
		{
			name: "go source, the shapes a Definition cannot describe",
			def:  goDefinition(),
			src:  "s := `raw\\n` // note\nr := '\\n'\nif s == r {\n\treturn 0x1F\n}",
			want: []Token{
				{Kind: KindIdent, Lexeme: "s"},
				{Kind: KindOperator, Lexeme: ":="},
				{Kind: KindString, Lexeme: `raw\n`},
				{Kind: KindIdent, Lexeme: "r"},
				{Kind: KindOperator, Lexeme: ":="},
				{Kind: KindString, Lexeme: `\n`},
				{Kind: KindKeyword, Lexeme: "if"},
				{Kind: KindIdent, Lexeme: "s"},
				{Kind: KindOperator, Lexeme: "=="},
				{Kind: KindIdent, Lexeme: "r"},
				{Kind: KindSymbol, Lexeme: "{"},
				{Kind: KindKeyword, Lexeme: "return"},
				{Kind: KindInt, Lexeme: "0"},
				{Kind: KindIdent, Lexeme: "x1F"},
				{Kind: KindSymbol, Lexeme: "}"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := lex(test.def, test.src)
			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

// goDefinition approximates the lexical surface of Go: both comment
// styles, the keywords, the operators with their shared prefixes,
// the punctuation, and — through Custom, the way the doc suggests —
// raw strings and rune literals, which a Definition cannot describe.
// It is a lexer for reading Go, not a front end: hex digits and
func goDefinition() Definition {
	return Definition{
		LineComment:    "//",
		BlockComment:   [2]string{"/*", "*/"},
		NestedComments: false,
		Keywords: []string{
			"break", "case", "chan", "const", "continue", "default",
			"defer", "else", "fallthrough", "for", "func", "go", "goto",
			"if", "import", "interface", "map", "package", "range",
			"return", "select", "struct", "switch", "type", "var",
		},
		Operators: []string{
			"<<=", ">>=", "&^=", "...", "&&", "||", "<-", "++", "--",
			"==", "!=", "<=", ">=", ":=", "+=", "-=", "*=", "/=",
			"%=", "&=", "|=", "^=", "<<", ">>", "&^", "+", "-", "*",
			"/", "%", "&", "|", "^", "<", ">", "=", "!", "~",
		},
		Punctuation: "()[]{}.,;:",
		Integers:    true,
		Strings:     true,
		Escapes: map[rune]rune{
			'n': '\n', 't': '\t', 'r': '\r', '\\': '\\', '"': '"',
			'\'': '\'', 'a': '\a', 'b': '\b', 'f': '\f', 'v': '\v',
			'0': 0, 'u': 'u', 'U': 'U', 'x': 'x',
		},
		Custom: []func(parsec.Buffer[rune, strings.Position]) (Token, error){rawString, runeLiteral},
	}
}

// rawString lexes a Go raw string literal: a backquote run whose
func rawString(buf parsec.Buffer[rune, strings.Position]) (Token, error) {
	r, err := buf.Read(true)
	if err != nil || r != '`' {
		return Token{}, ErrNoMatch
	}

	var text []rune

	for {
		if buf.IsEOF() {
			return Token{}, parsec.NewParseError(buf.Position(), "unterminated raw string")
		}

		r, err = buf.Read(true)
		if err != nil {
			return Token{}, err
		}

		if r == '`' {
			break
		}

		text = append(text, r)
	}

	return Token{Kind: KindString, Lexeme: Lexeme(text)}, nil
}

func runeLiteral(buf parsec.Buffer[rune, strings.Position]) (Token, error) {
	pos := buf.Position()

	r, err := buf.Read(true)
	if err != nil || r != '\'' {
		return Token{}, ErrNoMatch
	}

	var text []rune

	for {
		r, err = buf.Read(true)
		if err != nil {
			return Token{}, parsec.NewParseError(pos, "unterminated rune literal")
		}

		if r == '\'' {
			break
		}

		text = append(text, r)

		if r == '\\' {
			e, eerr := buf.Read(true)
			if eerr != nil {
				return Token{}, parsec.NewParseError(pos, "unterminated rune literal")
			}

			text = append(text, e)
		}
	}

	if len(text) == 0 {
		return Token{}, parsec.NewParseError(pos, "empty rune literal")
	}

	return Token{Kind: KindString, Lexeme: Lexeme(text)}, nil
}

// TestLexerGoSource reads a real Go file and checks the invariants:
// it lexes without errors, opens with the package keyword, and every
// keyword token is a reserved word of Go.
func TestLexerGoSource(t *testing.T) {
	src, err := os.ReadFile("lexer.go")
	require.NoError(t, err)

	def := goDefinition()
	stream, err := lex(def, string(src))
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(stream), 100)
	require.Equal(t, Token{Kind: KindKeyword, Lexeme: "package"}, stream[0])

	keywords := map[Lexeme]struct{}{}
	for _, kw := range def.Keywords {
		keywords[Lexeme(kw)] = struct{}{}
	}

	for _, tok := range stream {
		if tok.Kind != KindKeyword {
			continue
		}

		require.Contains(t, keywords, tok.Lexeme, "token %v is a keyword of no Go", tok)
	}
}
