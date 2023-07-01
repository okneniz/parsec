package parsec

// import (
	// "fmt"
// )

func Optional[T any, S any](c Combinator[T, S], def S) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		result, err := c(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

func Many[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			x, err := c(buffer)
			if err != nil {
				break
			}

			result = append(result, x)
		}

		return result, nil
	}
}

func Some[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		cc := Many(cap, c)

		// ignore err for coverage - many return at least empty slice
		result, _ := cc(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func Count[T any, S any](x int, next Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, x)

		for i := 0; i < x; i++ {
			n, err := next(buffer)
			if err != nil {
				return nil, err
			}

			result = append(result, n)
		}

		return result, nil
	}
}
