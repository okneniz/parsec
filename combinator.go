package parsec

type Combinator[T any, S any] func(Buffer[T]) (S, error)

type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T, S) B
type Composer3[T, S, B, M any] func(T, S, B) M

func Anything[T any](x T) bool { return true }
func Nothing[T any](x T) bool  { return false }

func Satisfy[T any](greedy bool, f Condition[T]) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, error) {
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

func Any[T any](greedy bool) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, error) {
		token, err := buffer.Read(greedy)
		if err != nil {
			return *new(T), err
		}

		return token, nil
	}
}

func Try[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		pos := buffer.Position()

		result, err := c(buffer)
		if err != nil {
			buffer.Seek(pos)
			return *new(S), err
		}

		return result, nil
	}
}

func Before[T any, S any, B any, Z any](
	body Combinator[T, B],
	before Combinator[T, S],
	compose Composer[S, B, Z],
) Combinator[T, Z] {
	return And(before, body, compose)
}

func After[T any, S any, B any, Z any](
	body Combinator[T, B],
	after Combinator[T, S],
	compose Composer[B, S, Z],
) Combinator[T, Z] {
	return And(body, after, compose)
}

func Between[T any, S any, B any, M any, Z any](
	pre Combinator[T, S],
	c Combinator[T, B],
	suf Combinator[T, M],
	compose Composer3[S, B, M, Z],
) Combinator[T, Z] {
	return func(buffer Buffer[T]) (Z, error) {
		prefix, err := pre(buffer)
		if err != nil {
			return *new(Z), err
		}

		body, err := c(buffer)
		if err != nil {
			return *new(Z), err
		}

		suffix, err := suf(buffer)
		if err != nil {
			return *new(Z), err
		}

		return compose(prefix, body, suffix), nil
	}
}

func Skip[T any, S any, B any](
	skip Combinator[T, B],
	next Combinator[T, S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		_, err := skip(buffer)
		if err != nil {
			return *new(S), err
		}

		return next(buffer)
	}
}

func SkipAfter[T any, S any, B any](
	skip Combinator[T, B], // TODO : change order of params?
	pre Combinator[T, S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		result, err := pre(buffer)
		if err != nil {
			return *new(S), err
		}

		_, err = skip(buffer)
		if err != nil {
			return *new(S), err
		}

		return result, nil
	}
}

func EOF[T any]() Combinator[T, bool] {
	return func(buffer Buffer[T]) (bool, error) {
		if buffer.IsEOF() {
			return true, nil
		}

		return false, nil
	}
}

func Cast[T any, S any, B any](
	c Combinator[T, S],
	f func(S) (B, error),
) Combinator[T, B] {
	return func(buffer Buffer[T]) (B, error) {
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
