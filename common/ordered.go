package common

import (
	"golang.org/x/exp/constraints"
)

func Range[T constraints.Ordered, P any](from T, to T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return x >= from && x <= to
	})
}

func NotRange[T constraints.Ordered, P any](from T, to T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return x < from || x > to
	})
}

func Gt[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t > x
	})
}

func Gte[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t >= x
	})
}

func Lt[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t < x
	})
}

func Lte[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t <= x
	})
}
