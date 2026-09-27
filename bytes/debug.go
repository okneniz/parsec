package bytes

import (
	"github.com/okneniz/parsec"
)

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[O any, S any](
	l parsec.Logged,
	m string,
	c parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Trace[byte, int, O, S](l, m, c)
}
