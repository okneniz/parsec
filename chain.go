package parsec

func Chainl[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
	def S,
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		f := Chainl1(c, op)

		result, ok := f(buffer)
		if !ok {
			return def, true
		}

		return result, ok
	}
}

func Chainl1[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		x, ok := c(buffer)
		if !ok {
			return *new(S), false
		}

		rest := x

		for !buffer.IsEOF() {
			f, ok := op(buffer)
			if !ok {
				break
			}

			y, ok := c(buffer)
			if !ok {
				break
			}

			rest = f(rest, y)
		}

		return rest, true
	}
}

func Chainr[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
	def S,
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		f := Chainr1(c, op)

		result, ok := f(buffer)
		if !ok {
			return def, true
		}

		return result, ok
	}
}

func Chainr1[T any, S any](
	c Combinator[T, S],
	op Combinator[T, func(S, S) S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		x, ok := c(buffer)
		if !ok {
			return *new(S), false
		}

		chain := make([]S, 0)
		chainF := make([]func(S, S) S, 0)

		chain = append(chain, x)

		for !buffer.IsEOF() {
			f, ok := op(buffer)
			if !ok {
				break
			}

			y, ok := c(buffer)
			if !ok {
				break
			}

			chainF = append(chainF, f)
			chain = append(chain, y)
		}

		if len(chainF) == 0 {
			return x, true
		}

		for len(chain) > 1 {
			a, b := chain[len(chain)-1], chain[len(chain)-2]
			g := chainF[len(chainF)-1]

			chain = chain[:len(chain)-2]
			chainF = chainF[:len(chainF)-1]

			chain = append(chain, g(b, a))
		}

		return chain[0], true
	}
}

func SepBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		token, ok := body(buffer)
		if !ok {
			return result, true
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
			token, ok = c(buffer)
			if !ok {
				break
			}

			result = append(result, token)
		}

		return result, true
	}
}

func SepBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		c := SepBy(cap, body, sep)

		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, false
		}
		return result, true
	}
}

func EndBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		c := Try(SkipAfter(sep, body))

		for !buffer.IsEOF() {
			token, ok := c(buffer)
			if !ok {
				break
			}

			result = append(result, token)
		}

		return result, true
	}
}

func EndBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		c := EndBy(cap, body, sep)

		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, false
		}
		return result, true
	}
}

func SepEndBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		s := Try(sep)

		for !buffer.IsEOF() {
			token, ok := body(buffer)
			if !ok {
				break
			}

			result = append(result, token)

			_, ok = s(buffer)
			if !ok {
				break
			}
		}

		return result, true
	}
}

func SepEndBy1[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		c := SepEndBy(cap, body, sep)

		result, _ := c(buffer)
		if len(result) == 0 {
			return nil, false
		}
		return result, true
	}
}

func ManyTill[T any, S any, B any](
	cap int,
	c Combinator[T, S],
	end Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		if buffer.IsEOF() {
			return nil, false
		}

		result := make([]S, 0, cap)
		z := Try(end)

		for {
			_, ok := z(buffer)
			if ok {
				break
			}

			x, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, x)
		}

		return result, true
	}
}
