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
