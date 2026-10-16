package parsec

// Concat applies the cs combinators one by one and concatenates
// their slice results into a single slice.
func Concat[T any, P any, O any, S any](
	limit int,
	cs ...Combinator[T, P, []O, S],
) Combinator[T, P, []O, S] {
	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, limit)

		for _, c := range cs {
			t, err := c(state, buffer)
			if err != nil {
				return nil, err
			}

			result = append(result, t...)
		}

		return result, nil
	}
}

// Sequence applies the cs combinators one by one
// and collects their results into a slice.
// It fails as soon as any of them fails.
func Sequence[T any, P any, O any, S any](
	limit int,
	cs ...Combinator[T, P, O, S],
) Combinator[T, P, []O, S] {
	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, limit)

		for _, c := range cs {
			t, err := c(state, buffer)
			if err != nil {
				return nil, err
			}

			result = append(result, t)
		}

		return result, nil
	}
}

// Choice - searches for a combinator that works successfully on the input data.
// If one is not found, it returns an ParseError error.
//
// Alternatives are tried one by one without restoring the buffer position
// between attempts, while greedy combinators consume one item even on failure.
// Wrap each alternative in Try, otherwise a failed alternative
// eats input and breaks the following ones:
//
//	Choice(errMessage, Try(Eq(...)), Try(OneOf(...)))
func Choice[T any, P any, O any, S any](
	errMessage string,
	cs ...Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		previous := make([]Error[P], 0)
		pos := buffer.Position()

		for _, c := range cs {
			result, err := c(state, buffer)
			if err == nil {
				return result, err
			}

			previous = append(previous, err)
		}

		return null, NewParseError(pos, errMessage, previous...)
	}
}

// Skip parses skip, discards its result,
// then parses body and returns the result of body.
func Skip[T any, P any, A any, O any, S any](
	skip Combinator[T, P, A, S],
	next Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		_, err := skip(state, buffer)
		if err != nil {
			return null, err
		}

		return next(state, buffer)
	}
}

// SkipAfter parses body first, then parses skip and discards its result.
// It returns the result of body.
func SkipAfter[T any, P any, A any, O any, S any](
	skip Combinator[T, P, A, S],
	body Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		result, err := body(state, buffer)
		if err != nil {
			return null, err
		}

		_, err = skip(state, buffer)
		if err != nil {
			return null, err
		}

		return result, nil
	}
}

// SkipMany skips a sequence of items parsed by the skip combinator
// and then applies body. Unlike Many it allocates nothing for
// the skipped part.
func SkipMany[T any, P any, A any, O any, S Recoverable[S]](
	skip Combinator[T, P, A, S],
	body Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	skip = Try(skip)

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		for !buffer.IsEOF() {
			_, err := skip(state, buffer)
			if err != nil {
				break
			}
		}

		return body(state, buffer)
	}
}

// Padded skips a sequence of items parsed by the skip combinator
// before and after body.
func Padded[T any, P any, A any, O any, S Recoverable[S]](
	skip Combinator[T, P, A, S],
	body Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	skip = Try(skip)

	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		for !buffer.IsEOF() {
			_, err := skip(state, buffer)
			if err != nil {
				break
			}
		}

		result, err := body(state, buffer)
		if err != nil {
			return null, err
		}

		for !buffer.IsEOF() {
			_, err := skip(state, buffer)
			if err != nil {
				break
			}
		}

		return result, nil
	}
}

// SkipSequence applies the cs combinators one by one,
// requires all of them to succeed and discards all the results.
func SkipSequence[T, P, O, S any](combs ...Combinator[T, P, O, S]) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		for _, c := range combs {
			_, err := c(state, buffer)
			if err != nil {
				return null, err
			}
		}

		return null, nil
	}
}

// SkipSequenceOf reads len(data) items and requires each of them
// to be equal to the corresponding element of data; the results
// are discarded. It reads the buffer directly, without extra allocations.
func SkipSequenceOf[T comparable, P, O, S any](
	errMessage string,
	data ...T,
) Combinator[T, P, O, S] {
	var null O

	return func(_ S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()

		for _, x := range data {
			r, err := buffer.Read(true)
			if err != nil {
				return null, NewParseError(pos, err.Error())
			}
			if x != r {
				return null, NewParseError(pos, errMessage)
			}
		}

		return null, nil
	}
}
