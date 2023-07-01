package parsec

func Or[T any, S any](x, y Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		result, err := x(buffer)
		if err != nil {
			return y(buffer)
		}

		return result, nil
	}
}

func And[T any, S any, B any, M any](
	x Combinator[T, S],
	y Combinator[T, B],
	compose Composer[S, B, M],
) Combinator[T, M] {
	return func(buffer Buffer[T]) (M, error) {
		first, err := x(buffer)
		if err != nil {
			return *new(M), err
		}

		second, err := y(buffer)
		if err != nil {
			return *new(M), err
		}

		return compose(first, second), nil
	}
}
