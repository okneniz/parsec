package common

// Concat - use cs combinators to parse slices step by step,
// concatenate all result to one big slice and returns it.
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

// Sequence - reads input elements one by one using cs combinators.
// If any of them fails, it returns an error.
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

// Choice - searches for a combinator that works successfully on the input data.
// if one is not found, it returns an NothingMatched error.
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

// Skip - ignores the result of the first combinator
// and returns only the result of the second.
func Skip[T any, P any, S any, B any](
	skip Combinator[T, P, B],
	next Combinator[T, P, S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		_, err := skip(buffer)
		if err != nil {
			return *new(S), err
		}

		return next(buffer)
	}
}

// SkipAfter - ignores the result of the first combinator
// and returns only the result of the second.
// Use body combinator at first.
func SkipAfter[T any, P any, S any, B any](
	skip Combinator[T, P, B],
	body Combinator[T, P, S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		result, err := body(buffer)
		if err != nil {
			return *new(S), err // TODO : allocate without new (by var)?
		}

		_, err = skip(buffer)
		if err != nil {
			return *new(S), err
		}

		return result, nil
	}
}

// SkipMany - skip sequence of items parsed by first combinator before body combinator.
func SkipMany[T any, P any, S any, B any](
	skip Combinator[T, P, S],
	body Combinator[T, P, B],
) Combinator[T, P, B] {
	skip = Try(skip)

	return func(buffer Buffer[T, P]) (B, error) {
		for !buffer.IsEOF() {
			_, err := skip(buffer)
			if err != nil {
				break
			}
		}

		return body(buffer)
	}
}

// Padded - skip sequence of items parsed by first combinator
// before and after body combinator.
func Padded[T any, P any, S any, B any](
	skip Combinator[T, P, S],
	body Combinator[T, P, B],
) Combinator[T, P, B] {
	skip = Try(skip)

	return func(buffer Buffer[T, P]) (B, error) {
		for !buffer.IsEOF() {
			_, err := skip(buffer)
			if err != nil {
				break
			}
		}

		result, err := body(buffer)
		if err != nil {
			return *new(B), err
		}

		for !buffer.IsEOF() {
			_, err := skip(buffer)
			if err != nil {
				break
			}
		}

		return result, nil
	}
}

// SkipSequence - reads input elements one by one using `cs` combinators and ignore it.
func SkipSequence[T, P, S any](combs ...Combinator[T, P, S]) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		var result S

		for _, c := range combs {
			_, err := c(buffer)
			if err != nil {
				return result, err
			}
		}

		return result, nil
	}
}

// SkipSequenceOf - reads input elements which must be equal input data and ignore it.
func SkipSequenceOf[T comparable, P, S any](data ...T) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		var result S

		for _, x := range data {
			r, err := buffer.Read(true)
			if err != nil {
				return result, err
			}
			if x != r {
				return result, NothingMatched
			}
		}

		return result, nil
	}
}
