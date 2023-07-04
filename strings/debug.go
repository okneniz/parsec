package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Trace[T any](
	l p.Logged,
	m string,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Trace[rune, Position, T](l, m, c)
}
