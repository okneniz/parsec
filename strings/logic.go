package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any](
	x p.Combinator[rune, Position, T],
	y p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Or[rune, Position, T](x, y)
}

// And - use x and y combinators to consume input data.
// Apply them result to compose function and return result of it.
func And[S any, B any, M any](
	x p.Combinator[rune, Position, S],
	y p.Combinator[rune, Position, B],
	compose p.Composer[S, B, M],
) p.Combinator[rune, Position, M] {
	return p.And[rune, Position, S, B, M](x, y, compose)
}
