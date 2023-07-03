package parsec

import (
	"golang.org/x/exp/constraints"
)

func Range[T constraints.Ordered](from T, to T) Combinator[T, T] {
	return Satisfy[T](true, func(x T) bool {
		return x >= from && x <= to
	})
}

func NotRange[T constraints.Ordered](from T, to T) Combinator[T, T] {
	return Satisfy[T](true, func(x T) bool {
		return x < from || x > to
	})
}
