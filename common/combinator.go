package common

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

func Any[T any, P any]() Combinator[T, P, T] {
	return func(buffer Buffer[T, P]) (T, error) {
		token, err := buffer.Read(true)
		if err != nil {
			return *new(T), err
		}

		return token, nil
	}
}

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

func EOF[T any, P any]() Combinator[T, P, bool] {
	return func(buffer Buffer[T, P]) (bool, error) {
		if buffer.IsEOF() {
			return true, nil
		}

		return false, nil
	}
}

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

func Const[T any, P any, S any](value S) Combinator[T, P, S] {
	return func(_ Buffer[T, P]) (S, error) {
		return value, nil
	}
}

func Fail[T any, P any, S any](err error) Combinator[T, P, S] {
	var x S

	return func(_ Buffer[T, P]) (S, error) {
		return x, err
	}
}
