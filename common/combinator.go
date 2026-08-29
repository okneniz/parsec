package common

// Satisfy succeeds for any item for which the supplied function f returns true
// and returns that item.
// If greedy is true, buffer keeps position after reading - even when f fails,
// so the failed combinator consumes exactly one item.
// Wrap it in Try if you need the position to be restored on failure.
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
			return null, NewParseError(pos, errMessage)
		}

		if f(token) {
			return token, nil
		}

		return null, NewParseError(pos, errMessage)
	}
}

// Any reads and returns the next item from the buffer.
// It fails only at the end of the buffer.
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

// Try applies c and, when it fails, rewinds the buffer to the previous position.
// This is the backtracking primitive of the library: greedy combinators like
// Satisfy, Eq, Range and others consume input even when they fail,
// so every combinator which can fail inside an alternative branch
// (see Choice, Or) or in a loop (see Many and similar) must be wrapped
// in Try to make backtracking possible.
func Try[T any, P any, S any](c Combinator[T, P, S]) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		pos := buffer.Position()

		result, err := c(buffer)
		if err != nil {
			if seekErr := buffer.Seek(pos); seekErr != nil {
				return null, NewParseError(
					buffer.Position(),
					seekErr.Error(),
					err,
				)
			}

			return null, err
		}

		return result, nil
	}
}

// Between parses open, then body, then close,
// and returns only the result of body.
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

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF[T any, P any]() Combinator[T, P, bool] {
	return func(buffer Buffer[T, P]) (bool, Error[P]) {
		if buffer.IsEOF() {
			return true, nil
		}

		return false, nil
	}
}

// Cast parses the input with c and transforms the result with f.
// It fails with errMessage when f returns an error.
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
			return value, NewParseError(pos, castError.Error())
		}

		return value, nil
	}
}

// Const consumes nothing and returns the given value.
func Const[T any, P any, S any](value S) Combinator[T, P, S] {
	return func(_ Buffer[T, P]) (S, Error[P]) {
		return value, nil
	}
}

// Fail consumes nothing and always fails with the given message.
func Fail[T any, P any, S any](errMessage string) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		return null, NewParseError(buffer.Position(), errMessage)
	}
}
