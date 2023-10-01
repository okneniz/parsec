package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Trace - writes messages to the log about the state of the buffer before
// and after using the combinator, the result of the cobinator and its error.
func Trace[T any](
	l p.Logged,
	m string,
	c p.Combinator[rune, int, T],
) p.Combinator[rune, int, T] {
	return p.Trace[rune, int, T](l, m, c)
}
