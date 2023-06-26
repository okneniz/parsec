package parsec

func Optional[T any, S any](c Combinator[T, S], def S) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		t, ok := c(buffer)
		if ok {
			return t, ok
		}

		return def, true
	}
}

func Many[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			t, ok := c(buffer)
			if !ok {
				break
			}

			result = append(result, t)
		}

		return result, true
	}
}

func Some[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		cc := Many(cap, c)

		t, _ := cc(buffer)
		if len(t) > 0 {
			return t, true
		}

		return nil, false
	}
}

func Count[T any, S any](x int, next Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		tokens := make([]S, 0, x)

		for i := 0; i < x; i++ {
			token, ok := next(buffer)
			if !ok {
				return nil, false
			}

			tokens = append(tokens, token)
		}

		return tokens, true
	}
}
