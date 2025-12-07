package common

import "fmt"

// Optional - use c combinator to consume input data from buffer.
// If it failed, than return def value.
func Optional[T any, P any, S any](c Combinator[T, P, S], def S) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, Error[P]) {
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
	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
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
func Some[T any, P any, S any](
	cap int,
	errMessage string,
	c Combinator[T, P, S],
) Combinator[T, P, []S] {
	parse := Many(cap, c)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		pos := buffer.Position()

		// ignore err for coverage - many return at least empty slice
		result, _ := parse(buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// Count - try to read X item by c combinator.
// Stop on first error.
func Count[T any, P any, S any](
	cap int,
	errMessage string,
	c Combinator[T, P, S],
) Combinator[T, P, []S] {
	f, err := Quantifier(errMessage, cap, cap, c)
	if err != nil {
		panic(err)
	}

	return f
}

// Quantifier - consume at items by c combinator,
// more than or equal than second param 'from' but less than or equal 'to'.
// Stop on first error.
func Quantifier[T any, P any, S any](
	errMessage string,
	from, to int,
	c Combinator[T, P, S],
) (Combinator[T, P, []S], error) {
	if from > to {
		return nil, fmt.Errorf(
			"param 'from' must be less than param 'to', actual from=%d, to=%d",
			from,
			to,
		)
	}

	if from < 0 {
		return nil, fmt.Errorf(
			"param 'from' must be positive, actual from=%d, to=%d",
			from,
			to,
		)
	}

	return func(buf Buffer[T, P]) ([]S, Error[P]) {
		start := buf.Position()
		result := make([]S, 0, to-from)

		for i := 0; i < to; i++ {
			pos := buf.Position()

			n, err := c(buf)
			if err != nil {
				if len(result) >= from {
					if seekErr := buf.Seek(pos); seekErr != nil {
						prevErr := NewParseError(buf.Position(), seekErr.Error(), err)
						return nil, NewParseError(start, errMessage, prevErr)
					}

					return result, nil
				}

				if seekErr := buf.Seek(start); seekErr != nil {
					prevErr := NewParseError(buf.Position(), seekErr.Error(), err)
					return nil, NewParseError(start, errMessage, prevErr)
				}

				return nil, NewParseError(start, errMessage, err)
			}

			result = append(result, n)
		}

		return result, nil
	}, nil
}
