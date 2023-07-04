package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Optional[T any](
	c p.Combinator[byte, int, T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Optional[byte, int, T](c, def)
}

func Many[T any](
	cap int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Many[byte, int, T](cap, c)
}

func Some[T any](
	cap int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Some[byte, int, T](cap, c)
}

func Count[T any](
	n int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Count[byte, int, T](n, c)
}
