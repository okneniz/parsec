package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Satisfy(
	greedy bool,
	f p.Condition[byte],
) p.Combinator[byte, int, byte] {
	return p.Satisfy[byte, int](greedy, f)
}

func Any() p.Combinator[byte, int, byte] {
	return p.Any[byte, int]()
}

func Try[T any](c p.Combinator[byte, int, T]) p.Combinator[byte, int, T] {
	return p.Try[byte, int, T](c)
}

func Between[T any, S any, B any](
	pre p.Combinator[byte, int, T],
	c p.Combinator[byte, int, S],
	suf p.Combinator[byte, int, B],
) p.Combinator[byte, int, S] {
	return p.Between(pre, c, suf)
}

func Skip[T any, S any](
	skip p.Combinator[byte, int, S],
	body p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.Skip(skip, body)
}

func SkipAfter[T any, S any](
	skip p.Combinator[byte, int, S],
	body p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.SkipAfter(skip, body)
}

func Padded[T any, S any](
	skip p.Combinator[byte, int, S],
	body p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.Padded(skip, body)
}

func EOF() p.Combinator[byte, int, bool] {
	return p.EOF[byte, int]()
}

func Cast[T any, S any](
	c p.Combinator[byte, int, T],
	f func(T) (S, error),
) p.Combinator[byte, int, S] {
	return p.Cast(c, f)
}