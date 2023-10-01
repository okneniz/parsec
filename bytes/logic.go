package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any](
	x p.Combinator[byte, int, T],
	y p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.Or[byte, int, T](x, y)
}

// And - use x and y combinators to consume input bytes.
// Apply them result to compose function and return result of it.
func And[S any, B any, M any](
	x p.Combinator[byte, int, S],
	y p.Combinator[byte, int, B],
	compose p.Composer[S, B, M],
) p.Combinator[byte, int, M] {
	return p.And[byte, int, S, B, M](x, y, compose)
}
