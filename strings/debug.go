package strings

import (
	"github.com/okneniz/parsec/common"
)

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[T any](
	l common.Logged,
	m string,
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Trace[rune, Position, T](l, m, c)
}
