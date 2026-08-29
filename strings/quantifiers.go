package strings

import (
	"github.com/okneniz/parsec"
)

// Optional applies c and returns its result,
// or the def value when c fails. It never fails.
func Optional[T any](
	c parsec.Combinator[rune, Position, T],
	def T,
) parsec.Combinator[rune, Position, T] {
	return parsec.Optional[rune, Position, T](c, def)
}

// Many applies c as many times as possible and collects the results.
// It stops at the first error or at the end of the buffer and returns
// everything collected so far, possibly an empty slice.
// Wrap c in Try to stop without consuming input.
func Many[T any](
	cap int,
	c parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, []T] {
	return parsec.Many[rune, Position, T](cap, c)
}

// Some is like Many but requires at least one item:
// it fails with errMessage when nothing could be parsed.
func Some[T any](
	cap int,
	errMessage string,
	c parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, []T] {
	return parsec.Some[rune, Position, T](cap, errMessage, c)
}

// Count applies c exactly cap times and collects the results.
// It fails with errMessage as soon as any application fails.
func Count[T any](
	cap int,
	errMessage string,
	c parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, []T] {
	return parsec.Count[rune, Position, T](cap, errMessage, c)
}
