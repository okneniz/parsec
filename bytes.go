package parsec

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
