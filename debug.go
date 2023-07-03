package parsec

import (
	"fmt"
)

type Logged interface {
	Log(args ...any)
}

func Trace[T any, S any](l Logged, m string, c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		l.Log(m)
		l.Log("\tposition before", buffer.Position())

		result, err := c(buffer)
		l.Log("\tposition after", buffer.Position())
		if err != nil {
			l.Log("\tnot parsed ", m, result, err)
			return *new(S), err
		}

		l.Log("\tparsed", result, fmt.Sprintf("%T", result))
		return result, err

	}
}
