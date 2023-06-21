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
				return *new(S), false
			}

			y, ok := c(buffer)
			if !ok {
				return *new(S), false
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
				return *new(S), false
			}

			y, ok := c(buffer)
			if !ok {
				return *new(S), false
			}

			chainF = append(chainF, f)
			chain = append(chain, y)
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
