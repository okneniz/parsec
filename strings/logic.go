package strings

import (
	"github.com/okneniz/parsec"
)

// Or returns the result of the first combinator;
// when it fails, the second combinator is used.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[O any, S any](
	errMessage string,
	x parsec.Combinator[rune, Position, O, S],
	y parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Or[rune, Position, O, S](errMessage, x, y)
}

// And parses x, then y, and combines their results
// with the compose function.
func And[A any, B any, O any, S any](
	x parsec.Combinator[rune, Position, A, S],
	y parsec.Combinator[rune, Position, B, S],
	compose parsec.Composer[A, B, O],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.And[rune, Position, A, B, O, S](x, y, compose)
}
