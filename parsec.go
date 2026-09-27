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
// Every combinator also carries a state: one value, given at the entry
// point, that travels through the whole parse. The state is mutable
// by convention and combinators that change it restore their changes
// before returning an error, so a failed alternative leaves the state
// as it found it, the way Try leaves the buffer position. A parser
// that needs no state uses [Stateless] for it.
//
// See the examples directory for complete parsers built with
// this library: JSON, timestamps, credit card numbers, math
// expressions, MessagePack, PNG, a tiny ML dialect and a math
// language built from a lang.Definition.
package parsec

// Parse applies the c combinator to the buffer with the given state
// and returns its result. It is a small convenience wrapper:
// c(state, buffer) does the same.
func Parse[T any, P any, O any, S any](
	state S,
	buffer Buffer[T, P],
	c Combinator[T, P, O, S],
) (O, Error[P]) {
	result, err := c(state, buffer)
	if err != nil {
		return result, err
	}

	return result, nil
}
