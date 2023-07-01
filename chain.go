package parsec

func Chainl[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
	def S,
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		f := Chainl1(c, op)

		result, err := f(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

func Chainl1[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
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

func Chainr[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
	def S,
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
		f := Chainr1(c, op)

		result, err := f(buffer)
		if err != nil {
			return def, nil
		}

		return result, nil
	}
}

func Chainr1[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, error) {
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

func SepBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		token, err := body(buffer)
		if err != nil {
			return result, nil
		}
		result = append(result, token)

		c := Try(
			And(
				sep,
				body,
				func(_ B, x S) S { return x },
			),
		)

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

func SepBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		c := SepBy(cap, body, sep)

		// ignore error because SepBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func EndBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		c := Try(SkipAfter(sep, body))

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

func EndBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		c := EndBy(cap, body, sep)

		// ignore error because EndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}

		return result, nil
	}
}

func SepEndBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)

		s := Try(sep)

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

func SepEndBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		c := SepEndBy(cap, body, sep)

		// ignore error because SepEndBy return empty list anyway
		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, NotEnoughElements
		}
		return result, nil
	}
}

func ManyTill[T any, S any, B any](
	cap int,
	c Combinator[T, S],
	end Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, error) {
		result := make([]S, 0, cap)
		z := Try(end)

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
