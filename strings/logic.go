package strings

import (
	"github.com/okneniz/parsec/common"
)

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any](
	errMessage string,
	x common.Combinator[rune, Position, T],
	y common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Or[rune, Position, T](errMessage, x, y)
}

// And - use x and y combinators to consume input data.
// Apply them result to compose function and return result of it.
func And[S any, B any, M any](
	x common.Combinator[rune, Position, S],
	y common.Combinator[rune, Position, B],
	compose common.Composer[S, B, M],
) common.Combinator[rune, Position, M] {
	return common.And[rune, Position, S, B, M](x, y, compose)
}
