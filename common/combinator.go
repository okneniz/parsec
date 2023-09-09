package common

// Satisfy - succeeds for any item for which the supplied function f returns true.
// Returns the item that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy[T any, P any](greedy bool, f Condition[T]) Combinator[T, P, T] {
	return func(buffer Buffer[T, P]) (T, error) {
		token, err := buffer.Read(greedy)
		if err != nil {
			return *new(T), err
		}

		if f(token) {
			return token, nil
		}

		return *new(T), NothingMatched
	}
}

// Any - returns the readed item.
func Any[T any, P any]() Combinator[T, P, T] {
	return func(buffer Buffer[T, P]) (T, error) {
		token, err := buffer.Read(true)
		if err != nil {
			return *new(T), err
		}

		return token, nil
	}
}

// Try - try to use c combinator, if it falls, it returns buffer to the previous position.
func Try[T any, P any, S any](c Combinator[T, P, S]) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		pos := buffer.Position()

		result, err := c(buffer)
		if err != nil {
			buffer.Seek(pos)
			return *new(S), err
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
	return func(buffer Buffer[T, P]) (B, error) {
		_, err := pre(buffer)
		if err != nil {
			return *new(B), err
		}

		body, err := c(buffer)
		if err != nil {
			return *new(B), err
		}

		_, err = suf(buffer)
		if err != nil {
			return *new(B), err
		}

		return body, nil
	}
}

// EOF - checks that buffer reading has finished.
func EOF[T any, P any]() Combinator[T, P, bool] {
	return func(buffer Buffer[T, P]) (bool, error) {
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
	f func(S) (B, error),
) Combinator[T, P, B] {
	return func(buffer Buffer[T, P]) (B, error) {
		result, err := c(buffer)
		if err != nil {
			return *new(B), err
		}

		value, err := f(result)
		if err != nil {
			return *new(B), err
		}

		return value, nil
	}
}

// Const - doesn't read anything, just return the input value.
func Const[T any, P any, S any](value S) Combinator[T, P, S] {
	return func(_ Buffer[T, P]) (S, error) {
		return value, nil
	}
}

// Fail - doesn't read anything, just return input error.
func Fail[T any, P any, S any](err error) Combinator[T, P, S] {
	var x S

	return func(_ Buffer[T, P]) (S, error) {
		return x, err
	}
}
