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
