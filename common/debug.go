package common

import (
	"fmt"
)

type Logged interface {
	Log(args ...any)
}

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

		l.Log("\tparsed:", result, fmt.Sprintf("%T", result))
		return result, err

	}
}
