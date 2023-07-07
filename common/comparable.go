package common

func Eq[T comparable, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t == x
	})
}

func NotEq[T comparable, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t != x
	})
}

func OneOf[T comparable, P any](data ...T) Combinator[T, P, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P](true, func(x T) bool {
		_, exists := m[x]
		return exists
	})
}

func NoneOf[T comparable, P any](data ...T) Combinator[T, P, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P](true, func(x T) bool {
		_, exists := m[x]
		return !exists
	})
}

func SequenceOf[T comparable, P any](data ...T) Combinator[T, P, []T] {
	return func(buffer Buffer[T, P]) ([]T, error) {
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

func Map[T any, P any, K comparable, V any](
	cases map[K]V,
	c Combinator[T, P, K],
) Combinator[T, P, V] {
	return func(buffer Buffer[T, P]) (V, error) {
		token, err := c(buffer)
		if err != nil {
			return *new(V), err
		}

		result, exists := cases[token]
		if !exists {
			return *new(V), NothingMatched
		}

		return result, nil
	}
}
