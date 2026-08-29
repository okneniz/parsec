package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Or returns the result of the first combinator;
// when it fails, the second combinator is used.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[T any](
	errMessage string,
	x common.Combinator[byte, int, T],
	y common.Combinator[byte, int, T],
) common.Combinator[byte, int, T] {
	return common.Or[byte, int, T](errMessage, x, y)
}

// And parses x, then y, and combines their results
// with the compose function.
func And[S any, B any, M any](
	x common.Combinator[byte, int, S],
	y common.Combinator[byte, int, B],
	compose common.Composer[S, B, M],
) common.Combinator[byte, int, M] {
	return common.And[byte, int, S, B, M](x, y, compose)
}
