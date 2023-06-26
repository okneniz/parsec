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
