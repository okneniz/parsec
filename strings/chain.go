package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Chainl[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Chainl[rune, Position, T](c, op, def)
}

func Chainl1[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
) p.Combinator[rune, Position, T] {
	return p.Chainl1[rune, Position, T](c, op)
}

func Chainr[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Chainr[rune, Position, T](c, op, def)
}

func Chainr1[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
) p.Combinator[rune, Position, T] {
	return p.Chainr1[rune, Position, T](c, op)
}

func SepBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepBy[rune, Position, T](cap, body, sep)
}

func SepBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepBy1[rune, Position, T](cap, body, sep)
}

func EndBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.EndBy[rune, Position, T](cap, body, sep)
}

func EndBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.EndBy1[rune, Position, T](cap, body, sep)
}

func SepEndBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepEndBy[rune, Position, T](cap, body, sep)
}

func SepEndBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepEndBy1[rune, Position, T](cap, body, sep)
}

func ManyTill[T any, S any](
	cap int,
	c p.Combinator[rune, Position, T],
	end p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.ManyTill[rune, Position, T](cap, c, end)
}
