// Package parsec is a Golang parser combinator library
// inspired by Haskell parsec.
//
// A parser is a function that takes some input and produces some
// structured output; a combinator is a function that combines parsers
// into bigger parsers. This module splits the functionality into
// three packages:
//
//   - [github.com/okneniz/parsec/strings] - combinators for rune input,
//     with line/column positions and text helpers;
//   - [github.com/okneniz/parsec/bytes] - combinators for binary input,
//     with big/little-endian readers;
//   - [github.com/okneniz/parsec/common] - the generic core shared
//     by both of them.
//
// See the examples directory for complete parsers built with
// this library: JSON, timestamps, credit card numbers, math
// expressions, MessagePack and PNG.
package parsec
