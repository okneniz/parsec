package strings

import (
	"github.com/okneniz/parsec"
)

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[O any, S any](
	l parsec.Logged,
	m string,
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Trace[rune, Position, O, S](l, m, c)
}
