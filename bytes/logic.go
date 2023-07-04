package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Or[T any](
	x p.Combinator[byte, int, T],
	y p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.Or[byte, int, T](x, y)
}

func And[S any, B any, M any](
	x p.Combinator[byte, int, S],
	y p.Combinator[byte, int, B],
	compose p.Composer[S, B, M],
) p.Combinator[byte, int, M] {
	return p.And[byte, int, S, B, M](x, y, compose)
}
