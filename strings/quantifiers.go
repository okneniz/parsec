package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Optional[T any](
	c p.Combinator[rune, Position, T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Optional[rune, Position, T](c, def)
}

func Many[T any](
	cap int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Many[rune, Position, T](cap, c)
}

func Some[T any](
	cap int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Some[rune, Position, T](cap, c)
}

func Count[T any](
	n int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Count[rune, Position, T](n, c)
}
