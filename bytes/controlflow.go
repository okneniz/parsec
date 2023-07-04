package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Concat[T any](
	cap int,
	cs ...p.Combinator[byte, int, []T],
) p.Combinator[byte, int, []T] {
	return p.Concat[byte, int, T](cap, cs...)
}

func Sequence[T any](
	cap int,
	cs ...p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Sequence[byte, int, T](cap, cs...)
}

func Choice[T any](
	cs ...p.Combinator[byte, int, T],
) p.Combinator[byte, int, T] {
	return p.Choice[byte, int, T](cs...)
}
