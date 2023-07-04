package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Satisfy(
	greedy bool,
	f p.Condition[rune],
) p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](greedy, f)
}

func Any() p.Combinator[rune, Position, rune] {
	return p.Any[rune, Position]()
}

func Try[T any](c p.Combinator[rune, Position, T]) p.Combinator[rune, Position, T] {
	return p.Try[rune, Position, T](c)
}

func Between[T any, S any, B any](
	pre p.Combinator[rune, Position, T],
	c p.Combinator[rune, Position, S],
	suf p.Combinator[rune, Position, B],
) p.Combinator[rune, Position, S] {
	return p.Between(pre, c, suf)
}

func Skip[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Skip(skip, body)
}

func SkipAfter[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.SkipAfter(skip, body)
}

func Padded[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Padded(skip, body)
}

func EOF() p.Combinator[rune, Position, bool] {
	return p.EOF[rune, Position]()
}

func Cast[T any, S any](
	c p.Combinator[rune, Position, T],
	f func(T) (S, error),
) p.Combinator[rune, Position, S] {
	return p.Cast(c, f)
}
