package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Chainl[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Chainl[byte, int, T](c, op, def)
}

func Chainl1[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
) p.Combinator[byte, int, T] {
	return p.Chainl1[byte, int, T](c, op)
}

func Chainr[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Chainr[byte, int, T](c, op, def)
}

func Chainr1[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
) p.Combinator[byte, int, T] {
	return p.Chainr1[byte, int, T](c, op)
}

func SepBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepBy[byte, int, T](cap, body, sep)
}

func SepBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepBy1[byte, int, T](cap, body, sep)
}

func EndBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.EndBy[byte, int, T](cap, body, sep)
}

func EndBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.EndBy1[byte, int, T](cap, body, sep)
}

func SepEndBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepEndBy[byte, int, T](cap, body, sep)
}

func SepEndBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepEndBy1[byte, int, T](cap, body, sep)
}

func ManyTill[T any, S any](
	cap int,
	c p.Combinator[byte, int, T],
	end p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.ManyTill[byte, int, T](cap, c, end)
}
