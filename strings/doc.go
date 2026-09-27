// Package strings provides parser combinators for rune input:
// a rune Buffer with line/column positions, text-specific helpers
// such as Parens or Unsigned, unicode character classes, and thin
// typed wrappers around the generic combinators of
// [github.com/okneniz/parsec].
//
// Parsing starts from [ParseString] (or [Parse]), a combinator and
// the parser state, parsec.Stateless when the parser keeps no
// state:
//
//	parser := strings.Padded(
//		strings.Try(strings.Space[parsec.Stateless]("whitespace")),
//		strings.Unsigned[int, parsec.Stateless](),
//	)
//
//	result, err := strings.ParseString(parsec.Stateless{}, " 42 ", parser)
//	// result == 42
package strings
