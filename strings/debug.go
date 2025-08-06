package strings

import (
	"github.com/okneniz/parsec/common"
)

// Trace - writes messages to the log about the state of the buffer before
// and after using the combinator, the result of the cobinator and its error.
func Trace[T any](
	l common.Logged,
	m string,
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Trace[rune, Position, T](l, m, c)
}
