package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any](
	errMessage string,
	x common.Combinator[byte, int, T],
	y common.Combinator[byte, int, T],
) common.Combinator[byte, int, T] {
	return common.Or[byte, int, T](errMessage, x, y)
}

// And - use x and y combinators to consume input bytes.
// Apply them result to compose function and return result of it.
func And[S any, B any, M any](
	x common.Combinator[byte, int, S],
	y common.Combinator[byte, int, B],
	compose common.Composer[S, B, M],
) common.Combinator[byte, int, M] {
	return common.And[byte, int, S, B, M](x, y, compose)
}
