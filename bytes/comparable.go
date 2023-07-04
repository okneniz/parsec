package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Eq(t byte) p.Combinator[byte, int, byte] {
	return p.Eq[byte, int](t)
}

func NotEq(t byte) p.Combinator[byte, int, byte] {
	return p.NotEq[byte, int](t)
}

func OneOf(data ...byte) p.Combinator[byte, int, byte] {
	return p.OneOf[byte, int](data...)
}

func NoneOf(data ...byte) p.Combinator[byte, int, byte] {
	return p.NoneOf[byte, int](data...)
}

func SequenceOf(data ...byte) p.Combinator[byte, int, []byte] {
	return p.SequenceOf[byte, int](data...)
}
