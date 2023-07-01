package parsec

type Logged interface {
	Log(args ...any)
}

func Trace[T any, S any](l Logged, m string, c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		l.Log(m)
		l.Log("\tposition", buffer.Position())

		result, err := c(buffer)
		if err != nil {
			l.Log("\tnot parsed", result, err)
			return *new(S), err
		}

		l.Log("\tparsed", result)
		return result, err

	}
}
