package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Range(from byte, to byte) p.Combinator[byte, int, byte] {
	return p.Range[byte, int](from, to)
}

func NotRange(from byte, to byte) p.Combinator[byte, int, byte] {
	return p.NotRange[byte, int](from, to)
}

func Gt(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x > t
	})
}

func Gte(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x >= t
	})
}

func Lt(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x < t
	})
}

func Lte(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x <= t
	})
}
