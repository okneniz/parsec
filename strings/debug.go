package strings

import (
	"github.com/okneniz/parsec"
)

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[T any](
	l parsec.Logged,
	m string,
	c parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Trace[rune, Position, T](l, m, c)
}
