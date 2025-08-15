package common

// Satisfy - succeeds for any item for which the supplied function f returns true.
// Returns the item that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy[T any, P any](
	errMessage string,
	greedy bool,
	f Condition[T],
) Combinator[T, P, T] {
	var null T

	return func(buffer Buffer[T, P]) (T, Error[P]) {
		pos := buffer.Position()

		token, err := buffer.Read(greedy)
		if err != nil {
			return *new(T), NewParseError(pos, errMessage)
		}

		if f(token) {
			return token, nil
		}

		return null, NewParseError(pos, errMessage)
	}
}

// Any - returns the readed item.
func Any[T any, P any]() Combinator[T, P, T] {
	var null T

	return func(buffer Buffer[T, P]) (T, Error[P]) {
		pos := buffer.Position()

		token, err := buffer.Read(true)
		if err != nil {
			return null, NewParseError(pos, err.Error())
		}

		return token, nil
	}
}

// Try - try to use c combinator, if it falls, it returns buffer to the previous position.
func Try[T any, P any, S any](c Combinator[T, P, S]) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		pos := buffer.Position()

		result, err := c(buffer)
		if err != nil {
			if seekErr := buffer.Seek(pos); seekErr != nil {
				return null, NewParseError(buffer.Position(), seekErr.Error())
			}

			return null, err
		}

		return result, nil
	}
}

// Between - parse sequence of input combinators, skip first and last results.
func Between[T any, P any, S any, B any, M any](
	pre Combinator[T, P, S],
	c Combinator[T, P, B],
	suf Combinator[T, P, M],
) Combinator[T, P, B] {
	var null B

	return func(buffer Buffer[T, P]) (B, Error[P]) {
		_, err := pre(buffer)
		if err != nil {
			return null, err
		}

		body, err := c(buffer)
		if err != nil {
			return null, err
		}

		_, err = suf(buffer)
		if err != nil {
			return null, err
		}

		return body, nil
	}
}

// EOF - checks that buffer reading has finished.
func EOF[T any, P any]() Combinator[T, P, bool] {
	return func(buffer Buffer[T, P]) (bool, Error[P]) {
		if buffer.IsEOF() {
			return true, nil
		}

		return false, nil
	}
}

// Cast - parse data by c combinator and apply to f function.
// Return result of f function.
func Cast[T any, P any, S any, B any](
	c Combinator[T, P, S],
	cast func(S) (B, error),
) Combinator[T, P, B] {
	var null B

	return func(buffer Buffer[T, P]) (B, Error[P]) {
		pos := buffer.Position()

		result, err := c(buffer)
		if err != nil {
			return null, err
		}

		value, castError := cast(result)
		if castError != nil {
			return null, NewParseError(pos, castError.Error())
		}

		return value, nil
	}
}

// Const - doesn't read anything, just return the input value.
func Const[T any, P any, S any](value S) Combinator[T, P, S] {
	return func(_ Buffer[T, P]) (S, Error[P]) {
		return value, nil
	}
}

// Fail - doesn't read anything, just return input error.
func Fail[T any, P any, S any](errMessage string) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		return null, NewParseError(buffer.Position(), errMessage)
	}
}
