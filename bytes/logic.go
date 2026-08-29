package bytes

import (
	"github.com/okneniz/parsec"
)

// Or returns the result of the first combinator;
// when it fails, the second combinator is used.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[T any](
	errMessage string,
	x parsec.Combinator[byte, int, T],
	y parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.Or[byte, int, T](errMessage, x, y)
}

// And parses x, then y, and combines their results
// with the compose function.
func And[S any, B any, M any](
	x parsec.Combinator[byte, int, S],
	y parsec.Combinator[byte, int, B],
	compose parsec.Composer[S, B, M],
) parsec.Combinator[byte, int, M] {
	return parsec.And[byte, int, S, B, M](x, y, compose)
}
