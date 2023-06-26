package parsec

func Concat[T any, S any](cap int, cs ...Combinator[T, []S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)
		x := 0

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				break
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

func Sequence[T any, S any](cap int, cs ...Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				break
			}

			result = append(result, t)
		}

		if len(result) != len(cs) {
			return nil, false
		}

		return result, true
	}
}
