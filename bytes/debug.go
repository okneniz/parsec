package bytes

import (
	"github.com/okneniz/parsec"
)

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[T any](
	l parsec.Logged,
	m string,
	c parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.Trace[byte, int, T](l, m, c)
}
