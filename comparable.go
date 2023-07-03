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

func SequenceOf[T comparable](data ...T) Combinator[T, []T] {
	return func(buffer Buffer[T]) ([]T, error) {
		result := make([]T, 0, len(data))

		for _, x := range data {
			token, err := buffer.Read(true)
			if err != nil {
				return nil, err
			}

			if x != token {
				return nil, NothingMatched
			}

			result = append(result, token)
		}

		return result, nil
	}
}
