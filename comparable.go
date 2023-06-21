package parsec

func Eq[T comparable](greedy bool, t T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return t == x
	})
}

func NotEq[T comparable](greedy bool, t T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return t != x
	})
}

func OneOf[T comparable](greedy bool, data ...T) Combinator[T, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T](greedy, func(x T) bool {
		_, exists := m[x]
		return exists
	})
}

func NoneOf[T comparable](greedy bool, data ...T) Combinator[T, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T](greedy, func(x T) bool {
		_, exists := m[x]
		return !exists
	})
}

func Concat[T comparable, S any](cap int, cs ...Combinator[T, []S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, 0)
		x := 0

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t...)
			x++
		}

		if x != len(cs) {
			return nil, false
		}

		return result, true
	}
}

func Sequence[T comparable, S any](cap int, cs ...Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, len(cs))

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t)
		}

		if len(result) != len(cs) {
			return nil, false
		}

		return result, true
	}
}
