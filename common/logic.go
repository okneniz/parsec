package common

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
func Or[T any, P any, S any](
	errMessage string,
	x, y Combinator[T, P, S],
) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		pos := buffer.Position()

		result, xErr := x(buffer)
		if xErr == nil {
			return result, nil
		}

		result, yErr := y(buffer)
		if yErr == nil {
			return result, nil
		}

		return null, NewParseError(pos, errMessage, xErr, yErr)
	}
}

// And - use x and y combinators to consume input data.
// Apply them result to compose function and return result of it.
func And[T any, P any, S any, B any, M any](
	x Combinator[T, P, S],
	y Combinator[T, P, B],
	compose Composer[S, B, M],
) Combinator[T, P, M] {
	var null M

	return func(buffer Buffer[T, P]) (M, Error[P]) {
		first, xErr := x(buffer)
		if xErr != nil {
			return null, xErr
		}

		second, yErr := y(buffer)
		if yErr != nil {
			return null, yErr
		}

		return compose(first, second), nil
	}
}
