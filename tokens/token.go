// Package tokens is the token layer of parsec: the generic Token
// type, the token combinators, and the comment machinery.
//
// A lexer is a one-token stepper over the rune buffer,
// func(parsec.Buffer[rune, strings.Position]) (Token, error): it
// skips the leading trivia itself and reports the exhausted input
// with parsec.ErrEndOfFile, the very sentinel the buffers use. The
// token combinators — Exact, OfKind, Satisfy — take that stepper as
// their first argument and run over the same rune buffer it reads
// from: nothing is materialized, the rune position doubles as the
// token position, and parsec.Try rewinds through the stepper. A
// language parser is then a combinator,
// parsec.Combinator[rune, strings.Position, S], built from the token
// combinators; the lang package assembles one from a Definition,
// which is where the mapping from runes to tokens is configured:
// nothing below that package maps runes to tokens.
//
// LineComment and CommentBody parse the comment forms the languages
// share: the opening of a comment stays a separate probe, so that a
// comment missing from the input stays distinguishable from one
// broken in its body.
//
// The package is deliberately kind-agnostic: the Token is generic
// over its kind and its lexeme — any comparable types, typically a
// language-defined enum for the kind and string for the lexeme. The
// lang package defines the conventional kinds of its generated
// lexers; the tiny ml example uses them through its Definition.
package tokens

import "fmt"

// Token is one lexical token; the Kind and the Lexeme are the
// language's own comparable types. Tokens carry no source
// coordinates: the buffer position — the rune buffer position in the
// lazy adapters — plays that role.
type Token[K comparable, L comparable] struct {
	Kind   K
	Lexeme L
}

// String returns a human-readable representation of the token.
func (t Token[K, L]) String() string {
	return fmt.Sprintf("%v %q", t.Kind, fmt.Sprint(t.Lexeme))
}
