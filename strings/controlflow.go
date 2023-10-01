package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Concat - use cs combinators to parse slices step by step,
// concatenate all result to one big slice and returns it.
func Concat[T any](
	cap int,
	cs ...p.Combinator[rune, Position, []T],
) p.Combinator[rune, Position, []T] {
	return p.Concat[rune, Position, T](cap, cs...)
}

// Sequence - reads input elements one by one using cs combinators.
// If any of them fails, it returns an error.
func Sequence[T any](
	cap int,
	cs ...p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Sequence[rune, Position, T](cap, cs...)
}

// Choice - searches for a combinator that works successfully on the input data.
// if one is not found, it returns an NothingMatched error.
func Choice[T any](
	cs ...p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Choice[rune, Position, T](cs...)
}

// Skip - ignores the result of the first combinator
// and returns only the result of the second.
func Skip[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Skip(skip, body)
}

// SkipAfter - ignores the result of the first combinator
// and returns only the result of the second.
// Use body combinator at first.
func SkipAfter[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.SkipAfter(skip, body)
}

// SkipMany - skip sequence of items parsed by first combinator before body combinator.
func SkipMany[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.SkipMany(skip, body)
}

// Padded - skip sequence of items parsed by first combinator
// before and after body combinator.
func Padded[T any, S any](
	skip p.Combinator[rune, Position, S],
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return p.Padded(skip, body)
}
