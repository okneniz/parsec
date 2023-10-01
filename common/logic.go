package common

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any, P any, S any](x, y Combinator[T, P, S]) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		result, err := x(buffer)
		if err != nil {
			return y(buffer)
		}

		return result, nil
	}
}

// And - use x and y combinators to consume input data.
// Apply them result to compose function and return result of it.
func And[T any, P any, S any, B any, M any](
	x Combinator[T, P, S],
	y Combinator[T, P, B],
	compose Composer[S, B, M],
) Combinator[T, P, M] {
	return func(buffer Buffer[T, P]) (M, error) {
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
