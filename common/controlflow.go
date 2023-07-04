package common

func Concat[T any, P any, S any](cap int, cs ...Combinator[T, P, []S]) Combinator[T, P, []S] {
	return func(buffer Buffer[T, P]) ([]S, error) {
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

func Sequence[T any, P any, S any](cap int, cs ...Combinator[T, P, S]) Combinator[T, P, []S] {
	return func(buffer Buffer[T, P]) ([]S, error) {
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

func Choice[T any, P any, S any](cs ...Combinator[T, P, S]) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		for _, c := range cs {
			result, err := c(buffer)
			if err == nil {
				return result, err
			}
		}

		return *new(S), NothingMatched
	}
}
