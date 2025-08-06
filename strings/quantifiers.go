package strings

import (
	"github.com/okneniz/parsec/common"
)

// Optional - use c combinator to consume input data from buffer.
// If it failed, than return def value.
func Optional[T any](
	c common.Combinator[rune, Position, T],
	def T,
) common.Combinator[rune, Position, T] {
	return common.Optional[rune, Position, T](c, def)
}

// Many - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an empty slice even if nothing could be parsed.
func Many[T any](
	cap int,
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, []T] {
	return common.Many[rune, Position, T](cap, c)
}

// Some - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an error if at least one element could not be read.
func Some[T any](
	cap int,
	errMessage string,
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, []T] {
	return common.Some[rune, Position, T](cap, errMessage, c)
}

// Count - try to read X item by c combinator.
// Stop on first error.
func Count[T any](
	cap int,
	errMessage string,
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, []T] {
	return common.Count[rune, Position, T](cap, errMessage, c)
}
