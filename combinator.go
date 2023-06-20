package parsec

import (
	"golang.org/x/exp/constraints"
)

type Combinator[T any, S any] func(Buffer[T]) (S, bool)

type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T, S) B
type Composer3[T, S, B, M any] func(T, S, B) M

// func Nothing[T any](x T) bool { return false }
func Anything[T any](x T) bool { return true }

func first[T any](x, _ T) T { return x }

func Satisfy[T any](greedy bool, f Condition[T]) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		token, ok := buffer.Read(greedy)
		if !ok {
			return *new(T), false
		}

		if f(token) {
			return token, true
		}

		return *new(T), false
	}
}

func Any[T any](greedy bool) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		token, ok := buffer.Read(greedy)
		if !ok {
			return *new(T), false
		}

		return token, true
	}
}

func Try[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		pos := buffer.Position()

		r, ok := c(buffer)
		if !ok {
			buffer.Seek(pos)
		}

		return r, ok
	}
}

func Or[T any, S any](x, y Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		r, ok := x(buffer)
		if !ok {
			return y(buffer)
		}

		return r, ok
	}
}

func And[T any, S any, B any, M any](
	x Combinator[T, S],
	y Combinator[T, B],
	compose Composer[S, B, M],
) Combinator[T, M] {
	return func(buffer Buffer[T]) (M, bool) {
		first, ok := x(buffer)
		if !ok {
			return *new(M), false
		}

		second, ok := y(buffer)
		if !ok {
			return *new(M), false
		}

		return compose(first, second), true
	}
}

func Eq[T comparable](greedy bool, t T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return t == x
	})
}

func NotEq[T comparable](greedy bool, t T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return t != x
	})
}

func Range[T constraints.Ordered](greedy bool, from T, to T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return x >= from && x <= to
	})
}

func NotRange[T constraints.Ordered](greedy bool, from T, to T) Combinator[T, T] {
	return Satisfy[T](greedy, func(x T) bool {
		return x < from || x > to
	})
}

// TODO : rename to Sequence
func Slice[T comparable, S any](cs ...Combinator[T, S]) Combinator[T, []S] {
	// TODO : add cap param
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, len(cs))

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t)
		}

		// TODO : check length

		return result, true
	}
}

func Concat[T comparable, S any](cap int, cs ...Combinator[T, []S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, 0)

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t...)
		}

		// TODO : check length

		return result, true
	}
}

func Optional[T any, S any](c Combinator[T, S], def S) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		t, ok := c(buffer)
		if ok {
			return t, ok
		}

		return def, true
	}
}

func Many[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			t, ok := c(buffer)
			if !ok {
				break
			}

			result = append(result, t)
		}

		return result, true
	}
}

func Some[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		cc := Many(cap, c)

		t, _ := cc(buffer)
		if len(t) > 0 {
			return t, true
		}

		return nil, false
	}
}

func Before[T any, S any, B any, Z any](
	body Combinator[T, B],
	before Combinator[T, S],
	compose Composer[S, B, Z],
) Combinator[T, Z] {
	return And(before, body, compose)
}

func After[T any, S any, B any, Z any](
	body Combinator[T, B],
	after Combinator[T, S],
	compose Composer[B, S, Z],
) Combinator[T, Z] {
	return And(body, after, compose)
}

func Between[T any, S any, B any, M any, Z any](
	pre Combinator[T, S],
	c Combinator[T, B],
	suf Combinator[T, M],
	compose Composer3[S, B, M, Z],
) Combinator[T, Z] {
	return func(buffer Buffer[T]) (Z, bool) {
		prefix, ok := pre(buffer)
		if !ok {
			return *new(Z), false
		}

		body, ok := c(buffer)
		if !ok {
			return *new(Z), false
		}

		suffix, ok := suf(buffer)
		if !ok {
			return *new(Z), false
		}

		return compose(prefix, body, suffix), true
	}
}

func Count[T any, S any](x int, next Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		tokens := make([]S, 0, x)

		for i := 0; i < x; i++ {
			token, ok := next(buffer)
			if !ok {
				return nil, false
			}

			tokens = append(tokens, token)
		}

		return tokens, true
	}
}

func Skip[T any, S any, B any](
	skip Combinator[T, B],
	next Combinator[T, S],
) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		_, ok := skip(buffer)
		if ok {
			return next(buffer)
		}

		return *new(S), false
	}
}

func SepBy[T any, S any](cap int, body, sep Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		token, ok := body(buffer)
		if !ok {
			return result, true
		}
		result = append(result, token)

		c := Try[T, S](Skip(sep, body))

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

func SepBy1[T any, S any](cap int, body, sep Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		c := SepBy(cap, body, sep)

		result, ok := c(buffer)
		if !ok {
			return nil, false
		}
		if len(result) == 0 {
			return nil, false
		}
		return result, ok
	}
}

func EndBy[T any, S any](cap int, body, sep Combinator[T, S]) Combinator[T, []S] {
	return Many(cap, Try(And(body, sep, first[S])))
}

func EndBy1[T any, S any](cap int, body, sep Combinator[T, S]) Combinator[T, []S] {
	return Some(cap, Try(And(body, sep, first[S])))
}

func SepEndBy[T any, S any, B any](
	cap int,
	body Combinator[T, S],
	sep Combinator[T, B],
) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for !buffer.IsEOF() {
			token, ok := body(buffer)
			if !ok {
				break
			}

			result = append(result, token)

			_, ok = sep(buffer)
			if !ok {
				break
			}
		}

		return result, true
	}
}

func SepEndBy1[T any, S any](cap int, body, sep Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		c := SepEndBy(cap, body, sep)

		result, ok := c(buffer)
		if !ok {
			return nil, false
		}
		if len(result) == 0 {
			return nil, false
		}
		return result, ok
	}
}

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

func EOF[T any]() Combinator[T, struct{}] {
	return func(buffer Buffer[T]) (struct{}, bool) {
		x := buffer.IsEOF()
		return struct{}{}, x
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

func Cast[T any, S any, B any](
	c Combinator[T, S],
	f func(S) B,
) Combinator[T, B] {
	return func(buffer Buffer[T]) (B, bool) {
		result, ok := c(buffer)
		if !ok {
			return *new(B), false
		}

		return f(result), true
	}
}

func Choice[T any, S any](cs ...Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		for _, c := range cs {
			result, ok := c(buffer)
			if ok {
				return result, ok
			}
		}

		return *new(S), false
	}
}

type Logged interface {
	Log(args ...any)
}

func Trace[T any, S any](l Logged, m string, c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		l.Log(m)
		x, ok := buffer.Read(false)
		l.Log("\tposition", buffer.Position(), x, ok)

		result, ok := c(buffer)
		if ok {
			l.Log("\tparsed", result)
			return result, ok
		}

		l.Log("\tnot parsed", result)

		return *new(S), false
	}
}
