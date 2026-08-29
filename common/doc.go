// Package common contains the generic core of the parsec library:
// the [Combinator] function type, the [Buffer] input abstraction,
// parse errors with positions, and every combinator implemented once
// for an arbitrary token type.
//
// A parser combinator is a function that takes some input and produces
// some structured output; a combinator is a function that combines
// parsers into bigger parsers. Users normally do not work with this
// package directly: [github.com/okneniz/parsec/strings] and
// [github.com/okneniz/parsec/bytes] provide ready-to-use rune and byte
// buffers together with thin typed wrappers around the combinators
// defined here.
package common
