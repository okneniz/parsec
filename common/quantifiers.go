package common

import "fmt"

// Optional applies c and returns its result,
// or the def value when c fails. It never fails.
func Optional[T any, P any, S any](c Combinator[T, P, S], def S) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, Error[P]) {
		result, err := c(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Many applies c as many times as possible and collects the results.
// It stops at the first error or at the end of the buffer and returns
// everything collected so far, possibly an empty slice.
// Wrap c in Try to stop without consuming input.
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

// Some is like Many but requires at least one item:
// it fails with errMessage when nothing could be parsed.
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

// Count applies c exactly cap times and collects the results.
// It fails with errMessage as soon as any application fails.
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

// Quantifier applies c at least from times but not more than to times.
// It stops at the first error after collecting from items, restoring
// the buffer position to the end of the last successful application.
// It returns a build error when from is greater than to or negative.
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
