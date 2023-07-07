package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Eq(t rune) p.Combinator[rune, Position, rune] {
	return p.Eq[rune, Position](t)
}

func NotEq(t rune) p.Combinator[rune, Position, rune] {
	return p.NotEq[rune, Position](t)
}

func OneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.OneOf[rune, Position](data...)
}

func NoneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.NoneOf[rune, Position](data...)
}

func SequenceOf(data ...rune) p.Combinator[rune, Position, []rune] {
	return p.SequenceOf[rune, Position](data...)
}

func Map[K comparable, V any](
	cases map[K]V,
	c p.Combinator[rune, Position, K],
) p.Combinator[rune, Position, V] {
	return p.Map[rune, Position, K, V](cases, c)
}
