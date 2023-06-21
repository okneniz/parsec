package parsec

func Or[T any, S any](x, y Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		r, ok := x(buffer)
		if !ok {
			return y(buffer)
		}

		return r, ok
	}
}

func And[T any, S any, B any, M any](
	x Combinator[T, S],
	y Combinator[T, B],
	compose Composer[S, B, M],
) Combinator[T, M] {
	return func(buffer Buffer[T]) (M, bool) {
		first, ok := x(buffer)
		if !ok {
			return *new(M), false
		}

		second, ok := y(buffer)
		if !ok {
			return *new(M), false
		}

		return compose(first, second), true
	}
}
