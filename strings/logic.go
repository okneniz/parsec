package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Or[T any](
	x p.Combinator[rune, Position, T],
	y p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Or[rune, Position, T](x, y)
}

func And[S any, B any, M any](
	x p.Combinator[rune, Position, S],
	y p.Combinator[rune, Position, B],
	compose p.Composer[S, B, M],
) p.Combinator[rune, Position, M] {
	return p.And[rune, Position, S, B, M](x, y, compose)
}
