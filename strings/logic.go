package strings

import (
	"github.com/okneniz/parsec/common"
)

// Or returns the result of the first combinator;
// when it fails, the second combinator is used.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[T any](
	errMessage string,
	x common.Combinator[rune, Position, T],
	y common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Or[rune, Position, T](errMessage, x, y)
}

// And parses x, then y, and combines their results
// with the compose function.
func And[S any, B any, M any](
	x common.Combinator[rune, Position, S],
	y common.Combinator[rune, Position, B],
	compose common.Composer[S, B, M],
) common.Combinator[rune, Position, M] {
	return common.And[rune, Position, S, B, M](x, y, compose)
}
