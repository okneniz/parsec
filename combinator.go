package parsec

type Combinator[T any, S any] func(Buffer[T]) (S, bool)

type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T, S) B
type Composer3[T, S, B, M any] func(T, S, B) M

func Anything[T any](x T) bool { return true }
func Nothing[T any](x T) bool  { return false }

func Satisfy[T any](greedy bool, f Condition[T]) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		token, ok := buffer.Read(greedy)
		if !ok {
			return *new(T), false
		}

		if f(token) {
			return token, true
		}

		return *new(T), false
	}
}

func Any[T any](greedy bool) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		token, ok := buffer.Read(greedy)
		if !ok {
			return *new(T), false
		}

		return token, true
	}
}

func Try[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		pos := buffer.Position()

		r, ok := c(buffer)
		if !ok {
			buffer.Seek(pos)
		}

		return r, ok
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
	return func(buffer Buffer[T]) (Z, bool) {
		prefix, ok := pre(buffer)
		if !ok {
			return *new(Z), false
		}

		body, ok := c(buffer)
		if !ok {
			return *new(Z), false
		}

		suffix, ok := suf(buffer)
		if !ok {
			return *new(Z), false
		}

		return compose(prefix, body, suffix), true
	}
}

func Skip[T any, S any, B any](
	skip Combinator[T, B],
	next Combinator[T, S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		_, ok := skip(buffer)
		if ok {
			return next(buffer)
		}

		return *new(S), false
	}
}

func SkipAfter[T any, S any, B any](
	skip Combinator[T, B],
	next Combinator[T, S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		result, ok := next(buffer)
		if !ok {
			return *new(S), false
		}

		_, ok = skip(buffer)
		if !ok {
			return *new(S), false
		}

		return result, true
	}
}

func EOF[T any]() Combinator[T, struct{}] {
	return func(buffer Buffer[T]) (struct{}, bool) {
		x := buffer.IsEOF()
		return struct{}{}, x
	}
}

func Cast[T any, S any, B any](
	c Combinator[T, S],
	f func(S) B,
) Combinator[T, B] {
	return func(buffer Buffer[T]) (B, bool) {
		result, ok := c(buffer)
		if !ok {
			return *new(B), false
		}

		return f(result), true
	}
}

func Choice[T any, S any](cs ...Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		for _, c := range cs {
			result, ok := c(buffer)
			if ok {
				return result, ok
			}
		}

		return *new(S), false
	}
}
