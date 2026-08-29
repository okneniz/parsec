// Package parsec is a Golang parser combinator library
// inspired by Haskell parsec.
//
// A parser is a function that takes some input and produces some
// structured output; a combinator is a function that combines parsers
// into bigger parsers. This package is the generic core of the
// library: the [Combinator] function type, the [Buffer] input
// abstraction, parse errors with positions, and every combinator
// implemented once for an arbitrary item type. The satellite
// packages build on it:
//
//   - [github.com/okneniz/parsec/strings] - combinators for rune input,
//     with line/column positions and text helpers;
//   - [github.com/okneniz/parsec/bytes] - combinators for binary input,
//     with big/little-endian readers;
//   - [github.com/okneniz/parsec/lang] - lexer generation and
//     expression machinery for basic programming languages, with
//     escape hatches for everything a declarative definition cannot
//     describe;
//   - [github.com/okneniz/parsec/tokens] - the token-stream layer
//     shared by lang and by hand-written two-stage parsers: the Token
//     type, a buffer implementing [Buffer], and token-level
//     combinators.
//
// See the examples directory for complete parsers built with
// this library: JSON, timestamps, credit card numbers, math
// expressions, MessagePack, PNG, a tiny ML dialect and a math
// language built from a lang.Definition.
package parsec

// Parse applies the c combinator to the buffer and returns its result.
// It is a small convenience wrapper: c(buffer) does the same.
func Parse[T any, P any, S any](
	buffer Buffer[T, P],
	c Combinator[T, P, S],
) (S, Error[P]) {
	result, err := c(buffer)
	if err != nil {
		return result, err
	}

	return result, nil
}
