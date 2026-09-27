package lang

import (
	"errors"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"
)

// lex runs the Lexer stepper to the end of input: the test-side
// replacement of the deleted lang.Lex.
func lex[S any](def Definition[S], state S, src string) ([]Token, error) {
	buf := strings.Buffer([]rune(src))
	next := Lexer(def)

	stream := make([]Token, 0, 64)

	for {
		tok, err := next(state, buf)
		if errors.Is(err, parsec.ErrEndOfFile) {
			break
		}
		if err != nil {
			return nil, err
		}

		stream = append(stream, tok)
	}

	return stream, nil
}

// lexemes is a shorthand for the flat token text of a source.
func lexemes[S any](t *testing.T, def Definition[S], state S, src string) []string {
	t.Helper()

	stream, err := lex(def, state, src)
	require.NoError(t, err)

	out := make([]string, 0, len(stream))
	for _, tok := range stream {
		out = append(out, string(tok.Lexeme))
	}
	return out
}

func testDefinition() Definition[parsec.Stateless] {
	return Definition[parsec.Stateless]{
		LineComment:    "//",
		BlockComment:   [2]string{"(*", "*)"},
		NestedComments: true,

		Keywords:    []string{"let", "in"},
		Operators:   []string{"==", "<=", "+", "-", "*", "/"},
		Punctuation: "(),=",
		Integers:    true,
		Strings:     true,
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{"identifiers and keywords", "let ifx in_ x'", []string{"let", "ifx", "in_", "x'"}},
		{"longest operator match", "a == b <= c + d - e / f * g", []string{
			"a", "==", "b", "<=", "c", "+", "d", "-", "e", "/", "f", "*", "g",
		}},
		{"numbers", "1 007 42", []string{"1", "007", "42"}},
		{"strings are decoded", `"a\nb" "" "q\"q"`, []string{"a\nb", "", `q"q`}},
		{"punctuation", "(a,b) = c", []string{"(", "a", ",", "b", ")", "=", "c"}},
		{"line comments", "1 // rest is ignored\n2 //more", []string{"1", "2"}},
		{"block comments", "1 (* one *) 2", []string{"1", "2"}},
		{"nested comments", "1 (* a (* b (* c *) d *) e *) 2", []string{"1", "2"}},
		{"comment between operators", "a(*x*)+b", []string{"a", "+", "b"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, lexemes(t, testDefinition(), parsec.Stateless{}, test.src))
		})
	}
}

func TestLexErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"unterminated block", "1 (* oops", "unterminated comment"},
		{"unterminated string", `"abc`, "unterminated string"},
		{"unknown escape", `"a\q"`, "unknown escape"},
		{"stray operator rune", "a < b", "unrecognized character"},
		{"unrecognized character", "a @ b", "unrecognized character"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := lex(testDefinition(), parsec.Stateless{}, test.src)
			require.ErrorContains(t, err, test.want)
		})
	}
}

func charLiteral(state parsec.Stateless, buf parsec.Buffer[rune, strings.Position]) (Token, error) {
	pos := buf.Position()

	hash := strings.Eq[parsec.Stateless]("expected '#'", '#')
	if _, err := hash(state, buf); err != nil {
		return Token{}, ErrNoMatch
	}

	quote := strings.Eq[parsec.Stateless](`expected '"'`, '"')
	if _, err := quote(state, buf); err != nil {
		return Token{}, parsec.NewParseError(pos, "expected character literal")
	}

	if buf.IsEOF() {
		return Token{}, parsec.NewParseError(pos, "unterminated character literal")
	}

	r, err := strings.Any[parsec.Stateless]()(state, buf)
	if err != nil {
		return Token{}, err
	}

	if r == '"' {
		return Token{}, parsec.NewParseError(pos, "empty character literal")
	}

	if _, err := quote(state, buf); err != nil {
		return Token{}, parsec.NewParseError(pos, "unterminated character literal")
	}

	return Token{Kind: "char", Lexeme: Lexeme(r)}, nil
}

func TestCustomEscapes(t *testing.T) {
	def := testDefinition()
	def.Escapes = map[rune]rune{
		'x':  '!',
		'\\': '\\',
		'"':  '"',
	}

	tests := []struct {
		src  string
		want string
	}{
		{`"a\xb"`, "a!b"},
		{`"\\q\\p"`, `\q\p`},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			got, err := lex(def, parsec.Stateless{}, test.src)
			require.NoError(t, err)
			require.Equal(t, test.want, string(got[0].Lexeme))
		})
	}

	// an escape outside the table is an error
	_, err := lex(def, parsec.Stateless{}, `"a\nb"`)
	require.ErrorContains(t, err, "unknown escape")
}

func TestCustomTokens(t *testing.T) {
	def := testDefinition()
	def.Custom = []func(state parsec.Stateless, buf parsec.Buffer[rune, strings.Position]) (Token, error){charLiteral}

	got := lexemes(t, def, parsec.Stateless{}, `x #"a" == #"b"`)
	require.Equal(t, []string{"x", "a", "==", "b"}, got)

	// the # is still free for other uses when the custom declines
	require.Equal(t, []string{"x", "==", "y"}, lexemes(t, def, parsec.Stateless{}, "x == y"))

	// custom errors abort Lex
	_, err := lex(def, parsec.Stateless{}, `x #"ab"`)
	require.ErrorContains(t, err, "unterminated character literal")
}
