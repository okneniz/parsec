package parsec

// BinaryOp is a binary operation used by the Chain* combinators:
// op takes the previously accumulated value and the next parsed value.
type BinaryOp[T any] func(T, T) T

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[T any, P any, S any](
	def S,
	c Combinator[T, P, S],
	op Combinator[T, P, BinaryOp[S]],
) Combinator[T, P, S] {
	parse := Chainl1(c, op)

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		result, err := parse(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, BinaryOp[S]],
) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		x, err := c(buffer)
		if err != nil {
			return null, err
		}

		rest := x

		for !buffer.IsEOF() {
			f, err := op(buffer)
			if err != nil {
				break
			}

			y, err := c(buffer)
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
func Chainr[T any, P any, S any](
	def S,
	c Combinator[T, P, S],
	op Combinator[T, P, BinaryOp[S]],
) Combinator[T, P, S] {
	f := Chainr1(c, op)

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		result, err := f(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, BinaryOp[S]],
) Combinator[T, P, S] {
	var null S

	return func(buffer Buffer[T, P]) (S, Error[P]) {
		x, err := c(buffer)
		if err != nil {
			return null, err
		}

		chain := make([]S, 0)
		chainF := make([]func(S, S) S, 0)

		chain = append(chain, x)

		for !buffer.IsEOF() {
			f, err := op(buffer)
			if err != nil {
				break
			}

			y, err := c(buffer)
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
func SepBy[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := Try(
		And(
			sep,
			body,
			func(_ B, x S) S { return x },
		),
	)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		result := make([]S, 0, cap)

		token, err := body(buffer)
		if err != nil {
			return result, nil
		}
		result = append(result, token)

		for !buffer.IsEOF() {
			token, err = c(buffer)
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
func SepBy1[T any, P any, S any, B any](
	cap int,
	errMessage string,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	parse := SepBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		pos := buffer.Position()

		// ignore error because SepBy return empty list anyway
		result, _ := parse(buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := Try(SkipAfter(sep, body))

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			token, err := c(buffer)
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
func EndBy1[T any, P any, S any, B any](
	cap int,
	errMessage string,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := EndBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		pos := buffer.Position()

		// ignore error because EndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	s := Try(sep)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			token, err := body(buffer)
			if err != nil {
				break
			}

			result = append(result, token)

			_, err = s(buffer)
			if err != nil {
				break
			}
		}

		return result, nil
	}
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[T any, P any, S any, B any](
	cap int,
	errMessage string,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := SepEndBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		pos := buffer.Position()

		// ignore error because SepEndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[T any, P any, S any, B any](
	cap int,
	errMessage string,
	c Combinator[T, P, S],
	end Combinator[T, P, B],
) Combinator[T, P, []S] {
	needStop := Try(end)

	return func(buffer Buffer[T, P]) ([]S, Error[P]) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			_, err := needStop(buffer)
			if err == nil {
				break
			}

			pos := buffer.Position()

			data, err := c(buffer)
			if err != nil {
				return nil, NewParseError(pos, errMessage)
			}

			result = append(result, data)
		}

		return result, nil
	}
}
