package strings

import (
	"github.com/okneniz/parsec"
)

// Or returns the result of the first combinator;
// when it fails, the second combinator is used.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[T any](
	errMessage string,
	x parsec.Combinator[rune, Position, T],
	y parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Or[rune, Position, T](errMessage, x, y)
}

// And parses x, then y, and combines their results
// with the compose function.
func And[S any, B any, M any](
	x parsec.Combinator[rune, Position, S],
	y parsec.Combinator[rune, Position, B],
	compose parsec.Composer[S, B, M],
) parsec.Combinator[rune, Position, M] {
	return parsec.And[rune, Position, S, B, M](x, y, compose)
}
