package parsec

func Concat[T any, S any](cap int, cs ...Combinator[T, []S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		for _, c := range cs {
			t, err := c(buffer)
			if err != nil {
				return nil, err
			}

			result = append(result, t...)
		}

		return result, nil
	}
}

func Sequence[T any, S any](cap int, cs ...Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		for _, c := range cs {
			t, err := c(buffer)
			if err != nil {
				break
			}

			result = append(result, t)
		}

		if len(result) != len(cs) {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func Choice[T any, S any](cs ...Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		for _, c := range cs {
			result, err := c(buffer)
			if err == nil {
				return result, err
			}
		}

		return *new(S), NothingMatched
	}
}
