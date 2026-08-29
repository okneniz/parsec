package common

import (
	"cmp"
)

// Range succeeds when the next item is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the item even on failure (see Try).
func Range[T cmp.Ordered, P any](
	errMessage string,
	from, to T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x >= from && x <= to
	})
}

// NotRange succeeds when the next item is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the item even on failure (see Try).
func NotRange[T cmp.Ordered, P any](
	errMessage string,
	from, to T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x < from || x > to
	})
}

// Gt succeeds when the next item is greater than t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Gt[T cmp.Ordered, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x > t
	})
}

// Gte succeeds when the next item is greater than or equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Gte[T cmp.Ordered, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x >= t
	})
}

// Lt succeeds when the next item is less than t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Lt[T cmp.Ordered, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x < t
	})
}

// Lte succeeds when the next item is less than or equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Lte[T cmp.Ordered, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return x <= t
	})
}
