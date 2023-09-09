package common

import (
	"golang.org/x/exp/constraints"
)

// Range - succeeds for any item which include in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Range[T constraints.Ordered, P any](from T, to T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return x >= from && x <= to
	})
}

// NotRange - succeeds for any item which not included in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotRange[T constraints.Ordered, P any](from T, to T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return x < from || x > to
	})
}

// Gt - succeeds for any item which greater than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gt[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t > x
	})
}

// Gte - succeeds for any item which greater than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gte[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t >= x
	})
}

// Lt - succeeds for any item which less than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lt[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t < x
	})
}

// Lte - succeeds for any item which less than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lte[T constraints.Ordered, P any](t T) Combinator[T, P, T] {
	return Satisfy[T, P](true, func(x T) bool {
		return t <= x
	})
}
