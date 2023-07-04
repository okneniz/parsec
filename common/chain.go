package common

func Chainl[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, func(S, S) S],
	def S,
) Combinator[T, P, S] {
	f := Chainl1(c, op)

	return func(buffer Buffer[T, P]) (S, error) {
		result, err := f(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

func Chainl1[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, func(S, S) S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		x, err := c(buffer)
		if err != nil {
			return *new(S), err
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

func Chainr[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, func(S, S) S],
	def S,
) Combinator[T, P, S] {
	f := Chainr1(c, op)

	return func(buffer Buffer[T, P]) (S, error) {
		result, err := f(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

func Chainr1[T any, P any, S any](
	c Combinator[T, P, S],
	op Combinator[T, P, func(S, S) S],
) Combinator[T, P, S] {
	return func(buffer Buffer[T, P]) (S, error) {
		x, err := c(buffer)
		if err != nil {
			return *new(S), err
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

	return func(buffer Buffer[T, P]) ([]S, error) {
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

func SepBy1[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := SepBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, error) {
		// ignore error because SepBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func EndBy[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := Try(SkipAfter(sep, body))

	return func(buffer Buffer[T, P]) ([]S, error) {
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

func EndBy1[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := EndBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, error) {
		// ignore error because EndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func SepEndBy[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	s := Try(sep)

	return func(buffer Buffer[T, P]) ([]S, error) {
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

func SepEndBy1[T any, P any, S any, B any](
	cap int,
	body Combinator[T, P, S],
	sep Combinator[T, P, B],
) Combinator[T, P, []S] {
	c := SepEndBy(cap, body, sep)

	return func(buffer Buffer[T, P]) ([]S, error) {
		// ignore error because SepEndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}
		return result, nil
	}
}

func ManyTill[T any, P any, S any, B any](
	cap int,
	c Combinator[T, P, S],
	end Combinator[T, P, B],
) Combinator[T, P, []S] {
	z := Try(end)

	return func(buffer Buffer[T, P]) ([]S, error) {
		result := make([]S, 0, cap)

		for {
			_, err := z(buffer)
			if err == nil {
				break
			}

			x, err := c(buffer)
			if err != nil {
				break
			}

			result = append(result, x)
		}

		return result, nil
	}
}
