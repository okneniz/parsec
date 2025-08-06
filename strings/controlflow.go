package strings

import (
	"github.com/okneniz/parsec/common"
)

// Concat - use cs combinators to parse slices step by step,
// concatenate all result to one big slice and returns it.
func Concat[T any](
	cap int,
	cs ...common.Combinator[rune, Position, []T],
) common.Combinator[rune, Position, []T] {
	return common.Concat[rune, Position, T](cap, cs...)
}

// Sequence - reads input elements one by one using cs combinators.
// If any of them fails, it returns an error.
func Sequence[T any](
	cap int,
	cs ...common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, []T] {
	return common.Sequence[rune, Position, T](cap, cs...)
}

// Choice - searches for a combinator that works successfully on the input data.
// if one is not found, it returns an NothingMatched error.
func Choice[T any](
	cs ...common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	var null T

	return func(buffer common.Buffer[rune, Position]) (T, common.Error[Position]) {
		var deepestErr common.Error[Position]

		for _, c := range cs {
			result, err := c(buffer)
			if err == nil {
				return result, err
			}

			if deepestErr == nil {
				deepestErr = err
				continue
			}

			if err.Position().Deeper(deepestErr.Position()) {
				deepestErr = err
			}
		}

		return null, deepestErr
	}
}

// Skip - ignores the result of the first combinator
// and returns only the result of the second.
func Skip[T any, S any](
	skip common.Combinator[rune, Position, S],
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Skip(skip, body)
}

// SkipAfter - ignores the result of the first combinator
// and returns only the result of the second.
// Use body combinator at first.
func SkipAfter[T any, S any](
	skip common.Combinator[rune, Position, S],
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.SkipAfter(skip, body)
}

// SkipMany - skip sequence of items parsed by first combinator before body combinator.
func SkipMany[T any, S any](
	skip common.Combinator[rune, Position, S],
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.SkipMany(skip, body)
}

// Padded - skip sequence of items parsed by first combinator
// before and after body combinator.
func Padded[T any, S any](
	skip common.Combinator[rune, Position, S],
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Padded(skip, body)
}
