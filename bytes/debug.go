package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Trace[T any](
	l p.Logged,
	m string,
	c p.Combinator[rune, int, T],
) p.Combinator[rune, int, T] {
	return p.Trace[rune, int, T](l, m, c)
}
