package parsec

type Logged interface {
	Log(args ...any)
}

func Trace[T any, S any](l Logged, m string, c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		l.Log(m)
		x, ok := buffer.Read(false)
		l.Log("\tposition", buffer.Position(), x, ok)

		result, ok := c(buffer)
		if ok {
			l.Log("\tparsed", result)
			return result, ok
		}

		l.Log("\tnot parsed", result)

		return *new(S), false
	}
}
