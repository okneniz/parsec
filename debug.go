package parsec

import (
	"fmt"
)

// Logged is a minimal logging target accepted by Trace.
type Logged interface {
	Log(args ...any)
}

// Trace wraps c and logs the buffer position before and after its
// application, plus the parsed result or the error.
// It is a debugging helper and changes no parsing behavior.
func Trace[T any, P any, S any](
	l Logged,
	m string,
	c Combinator[T, P, S],
) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		l.Log(m)
		l.Log("\tposition before:", buffer.Position())

		result, err := c(buffer)
		l.Log("\tposition after:", buffer.Position())
		if err != nil {
			l.Log("\tnot parsed:", m, result, err)
			return null, err
		}

		l.Log("\tparsed:", fmt.Sprintf("%#v", result))
		return result, err
	}
}
