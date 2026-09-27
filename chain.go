package parsec

import "iter"

// BinaryOp is a binary operation used by the Chain* combinators:
// op takes the previously accumulated value and the next parsed value.
type BinaryOp[T any] func(T, T) T

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[T any, P any, O any, S any](
	def O,
	c Combinator[T, P, O, S],
	op Combinator[T, P, BinaryOp[O], S],
) Combinator[T, P, O, S] {
	parse := Chainl1(c, op)

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		result, err := parse(state, buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[T any, P any, O any, S any](
	c Combinator[T, P, O, S],
	op Combinator[T, P, BinaryOp[O], S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		x, err := c(state, buffer)
		if err != nil {
			return null, err
		}

		rest := x

		for !buffer.IsEOF() {
			f, err := op(state, buffer)
			if err != nil {
				break
			}

			y, err := c(state, buffer)
			if err != nil {
				break
			}

			rest = f(rest, y)
		}

		return rest, nil
	}
}

// Chainr parses zero or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// When nothing is parsed, def is returned.
func Chainr[T any, P any, O any, S any](
	def O,
	c Combinator[T, P, O, S],
	op Combinator[T, P, BinaryOp[O], S],
) Combinator[T, P, O, S] {
	f := Chainr1(c, op)

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		result, err := f(state, buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[T any, P any, O any, S any](
	c Combinator[T, P, O, S],
	op Combinator[T, P, BinaryOp[O], S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		x, err := c(state, buffer)
		if err != nil {
			return null, err
		}

		chain := make([]O, 0)
		chainF := make([]func(O, O) O, 0)

		chain = append(chain, x)

		for !buffer.IsEOF() {
			f, err := op(state, buffer)
			if err != nil {
				break
			}

			y, err := c(state, buffer)
			if err != nil {
				break
			}

			chainF = append(chainF, f)
			chain = append(chain, y)
		}

		if len(chainF) == 0 {
			return x, nil
		}

		for len(chain) > 1 {
			a, b := chain[len(chain)-1], chain[len(chain)-2]
			g := chainF[len(chainF)-1]

			chain = chain[:len(chain)-2]
			chainF = chainF[:len(chainF)-1]

			chain = append(chain, g(b, a))
		}

		return chain[0], nil
	}
}

// SepBy parses zero or more values of body separated by sep
// and returns them as a slice. A trailing separator is not allowed.
func SepBy[T any, P any, O any, A any, S any](
	cap int,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	c := Try(
		And(
			sep,
			body,
			func(_ A, x O) O { return x },
		),
	)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, cap)

		token, err := body(state, buffer)
		if err != nil {
			return result, nil
		}
		result = append(result, token)

		for !buffer.IsEOF() {
			token, err = c(state, buffer)
			if err != nil {
				break
			}

			result = append(result, token)
		}

		return result, nil
	}
}

// SepBy1 is like SepBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepBy1[T any, P any, O any, A any, S any](
	cap int,
	errMessage string,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	parse := SepBy(cap, body, sep)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		pos := buffer.Position()

		// ignore error because SepBy return empty list anyway
		result, _ := parse(state, buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[T any, P any, O any, A any, S any](
	cap int,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	c := Try(SkipAfter(sep, body))

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, cap)

		for !buffer.IsEOF() {
			token, err := c(state, buffer)
			if err != nil {
				break
			}

			result = append(result, token)
		}

		return result, nil
	}
}

// EndBy1 is like EndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func EndBy1[T any, P any, O any, A any, S any](
	cap int,
	errMessage string,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	c := EndBy(cap, body, sep)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		pos := buffer.Position()

		// ignore error because EndBy return empty list anyway
		result, _ := c(state, buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[T any, P any, O any, A any, S any](
	cap int,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	s := Try(sep)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, cap)

		for !buffer.IsEOF() {
			token, err := body(state, buffer)
			if err != nil {
				break
			}

			result = append(result, token)

			_, err = s(state, buffer)
			if err != nil {
				break
			}
		}

		return result, nil
	}
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[T any, P any, O any, A any, S any](
	cap int,
	errMessage string,
	body Combinator[T, P, O, S],
	sep Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	c := SepEndBy(cap, body, sep)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		pos := buffer.Position()

		// ignore error because SepEndBy return empty list anyway
		result, _ := c(state, buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[T any, P any, O any, A any, S any](
	cap int,
	errMessage string,
	c Combinator[T, P, O, S],
	end Combinator[T, P, A, S],
) Combinator[T, P, []O, S] {
	needStop := Try(end)

	return func(state S, buffer Buffer[T, P]) ([]O, Error[P]) {
		result := make([]O, 0, cap)

		for !buffer.IsEOF() {
			_, err := needStop(state, buffer)
			if err == nil {
				break
			}

			pos := buffer.Position()

			data, err := c(state, buffer)
			if err != nil {
				return nil, NewParseError(pos, errMessage)
			}

			result = append(result, data)
		}

		return result, nil
	}
}

// Seq applies c lazily and repeatedly: running the combinator on a
// buffer returns an iterator over the results of every step. Each
// successful step yields its value with a nil error; the first
// failing step is yielded with its error and ends the sequence, and
// the end of the buffer ends it silently. The iterator advances the
// shared buffer as the consumer ranges — an early break leaves the
// buffer right after the last yielded value, ready for the next
// combinator. Wrap c in Try to keep the buffer intact after the
// failing step. Like Many, it never ends when c succeeds without
// consuming input.
func Seq[T any, P any, O any, S any](c Combinator[T, P, O, S]) Combinator[T, P, iter.Seq2[O, Error[P]], S] {
	return func(state S, buffer Buffer[T, P]) (iter.Seq2[O, Error[P]], Error[P]) {
		seq := func(yield func(O, Error[P]) bool) {
			var null O

			for !buffer.IsEOF() {
				result, err := c(state, buffer)
				if err != nil {
					yield(null, err)

					return
				}

				if !yield(result, nil) {
					return
				}
			}
		}

		return seq, nil
	}
}
