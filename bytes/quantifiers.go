package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Optional applies c and returns its result,
// or the def value when c fails. It never fails.
func Optional[T any](
	c common.Combinator[byte, int, T],
	def T,
) common.Combinator[byte, int, T] {
	return common.Optional[byte, int, T](c, def)
}

// Many applies c as many times as possible and collects the results.
// It stops at the first error or at the end of the buffer and returns
// everything collected so far, possibly an empty slice.
// Wrap c in Try to stop without consuming input.
func Many[T any](
	cap int,
	c common.Combinator[byte, int, T],
) common.Combinator[byte, int, []T] {
	return common.Many[byte, int, T](cap, c)
}

// Some is like Many but requires at least one item:
// it fails with errMessage when nothing could be parsed.
func Some[T any](
	cap int,
	errMessage string,
	c common.Combinator[byte, int, T],
) common.Combinator[byte, int, []T] {
	return common.Some[byte, int, T](cap, errMessage, c)
}

// Count applies c exactly cap times and collects the results.
// It fails with errMessage as soon as any application fails.
func Count[T any](
	cap int,
	errMessage string,
	c common.Combinator[byte, int, T],
) common.Combinator[byte, int, []T] {
	return common.Count[byte, int, T](cap, errMessage, c)
}
