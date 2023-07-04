package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Concat[T any](
	cap int,
	cs ...p.Combinator[rune, Position, []T],
) p.Combinator[rune, Position, []T] {
	return p.Concat[rune, Position, T](cap, cs...)
}

func Sequence[T any](
	cap int,
	cs ...p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Sequence[rune, Position, T](cap, cs...)
}

func Choice[T any](
	cs ...p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Choice[rune, Position, T](cs...)
}
