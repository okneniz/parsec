package parsec

import (
	"golang.org/x/exp/constraints"
)

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
