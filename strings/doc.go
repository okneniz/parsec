// Package strings provides parser combinators for rune input:
// a rune Buffer with line/column positions, text-specific helpers
// such as Parens or Unsigned, unicode character classes, and thin
// typed wrappers around the generic combinators of
// [github.com/okneniz/parsec].
//
// Parsing starts from [ParseString] (or [Parse]) and a combinator:
//
//	parser := strings.Padded(
//		strings.Try(strings.Space("whitespace")),
//		strings.Unsigned[int](),
//	)
//
//	result, err := strings.ParseString(" 42 ", parser)
//	// result == 42
package strings
