package common

import (
	"fmt"
)

type Logged interface {
	Log(args ...any)
}

// Trace - writes messages to the log about the state of the buffer before
// and after using the combinator, the result of the cobinator and its error.
func Trace[T any, P any, S any](l Logged, m string, c Combinator[T, P, S]) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		l.Log(m)
		l.Log("\tposition before:", buffer.Position())

		result, err := c(buffer)
		l.Log("\tposition after:", buffer.Position())
		if err != nil {
			l.Log("\tnot parsed:", m, result, err)
			return *new(S), err
		}

		l.Log("\tparsed:", fmt.Sprintf("%#v", result))
		return result, err
	}
}
