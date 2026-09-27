package parsec

// Or - returns the result of the first combinator,
// if it fails, uses the second combinator.
// Like Choice, it does not restore the buffer position between attempts:
// wrap both arguments in Try if they can fail after consuming input.
func Or[T any, P any, O any, S any](
	errMessage string,
	x, y Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()

		result, xErr := x(state, buffer)
		if xErr == nil {
			return result, nil
		}

		result, yErr := y(state, buffer)
		if yErr == nil {
			return result, nil
		}

		return null, NewParseError(pos, errMessage, xErr, yErr)
	}
}

// And parses x, then y, and combines their results
// with the compose function.
func And[T any, P any, A any, B any, O any, S any](
	x Combinator[T, P, A, S],
	y Combinator[T, P, B, S],
	compose Composer[A, B, O],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		first, xErr := x(state, buffer)
		if xErr != nil {
			return null, xErr
		}

		second, yErr := y(state, buffer)
		if yErr != nil {
			return null, yErr
		}

		return compose(first, second), nil
	}
}
