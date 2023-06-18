package parsec

import (
	"fmt"
)

type Combinator[T any, S any] func(Buffer[T]) (S, bool)
type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T,S) B
type Composer3[T, S, B, M any] func(T, S, B) M

func Nothing[T any](x T) bool { return false }
func Anything[T any](x T) bool { return true }

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

func Not[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		r, ok := c(buffer)
		if ok {
			return *new(S), false
		}

		return r, true
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
	x Combinator[T,S],
	y Combinator[T,B],
	compose Composer[S,B,M],
) Combinator[T,M] {
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

func Token[T comparable](greedy bool, t T) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		r, ok := buffer.Read(greedy)
		if !ok || (r != t) {
			return *new(T), false
		}

		return r, true
	}
}

func Slice[T comparable, S any](cs ...Combinator[T,S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, len(cs))

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t)
		}

		return result, true
	}
}

func Optional[T any, S any](c Combinator[T, S]) Combinator[T, *S] {
	return func(buffer Buffer[T]) (*S, bool) {
		t, ok := c(buffer)
		if ok {
			return &t, ok
		}

		return nil, true
	}
}

func Many[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for {
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
	body Combinator[T,B],
	before Combinator[T,S],
	compose Composer[S,B,Z],
) Combinator[T, Z] {
	return And(before, body, compose)
}

func After[T any, S any, B any, Z any](
	body Combinator[T,B],
	after Combinator[T,S],
	compose Composer[B,S,Z],
) Combinator[T, Z] {
	return And(body, after, compose)
}

func Between[T any, S any, B any, M any, Z any](
	pre Combinator[T,S],
	c Combinator[T,B],
	suf Combinator[T,M],
	compose Composer3[S,B,M,Z],
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

		fmt.Println("WTF", body, ok)

		suffix, ok := suf(buffer)
		if !ok {
			return *new(Z), false
		}

		return compose(prefix, body, suffix), true
	}
}
