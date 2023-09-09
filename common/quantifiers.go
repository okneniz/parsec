package common

// Optional - use c combinator to consume input data from buffer.
// If it failed, than return def value.
func Optional[T any, P any, S any](c Combinator[T, P, S], def S) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		result, err := c(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Many - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an empty slice even if nothing could be parsed.
func Many[T any, P any, S any](cap int, c Combinator[T, P, S]) Combinator[T, P, []S] {
	return func(buffer Buffer[T, P]) ([]S, error) {
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

// Some - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an error if at least one element could not be read.
func Some[T any, P any, S any](cap int, c Combinator[T, P, S]) Combinator[T, P, []S] {
	return func(buffer Buffer[T, P]) ([]S, error) {
		cc := Many(cap, c)

		// ignore err for coverage - many return at least empty slice
		result, _ := cc(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

// Count - try to read X item by c combinator.
// Stop on first error.
func Count[T any, P any, S any](x int, c Combinator[T, P, S]) Combinator[T, P, []S] {
	return func(buffer Buffer[T, P]) ([]S, error) {
		result := make([]S, 0, x)

		for i := 0; i < x; i++ {
			n, err := c(buffer)
			if err != nil {
				return nil, err
			}

			result = append(result, n)
		}

		return result, nil
	}
}
