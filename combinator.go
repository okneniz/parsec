package parsec

type Combinator[T any, P any, S any] func(Buffer[T, P]) (S, error)

type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T, S) B

func Anything[T any](x T) bool { return true }
func Nothing[T any](x T) bool  { return false }

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

func Before[T any, P any, S any, B any, Z any](
	body Combinator[T, P, B], // TODO : change order of params to parsing order
	before Combinator[T, P, S],
	compose Composer[S, B, Z],
) Combinator[T, P, Z] {
	return And(before, body, compose)
}

func After[T any, P any, S any, B any, Z any](
	body Combinator[T, P, B],
	after Combinator[T, P, S],
	compose Composer[B, S, Z],
) Combinator[T, P, Z] {
	return And(body, after, compose)
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

func Skip[T any, P any, S any, B any](
	skip Combinator[T, P, B],
	next Combinator[T, P, S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		_, err := skip(buffer)
		if err != nil {
			return *new(S), err
		}

		return next(buffer)
	}
}

func SkipAfter[T any, P any, S any, B any](
	skip Combinator[T, P, B],
	body Combinator[T, P, S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		result, err := body(buffer)
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

func Padded[T any, P any, S any, B any](
	skip Combinator[T, P, B],
	body Combinator[T, P, S],
) Combinator[T, P, S] {
	x := Many(0, Try(skip))
	return Skip(x, SkipAfter(x, body))
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
