package strings

import (
	"github.com/okneniz/parsec"
)

// Concat applies the cs combinators one by one and concatenates
// their slice results into a single slice.
func Concat[T any](
	cap int,
	cs ...parsec.Combinator[rune, Position, []T],
) parsec.Combinator[rune, Position, []T] {
	return parsec.Concat[rune, Position, T](cap, cs...)
}

// Sequence applies the cs combinators one by one
// and collects their results into a slice.
// It fails as soon as any of them fails.
func Sequence[T any](
	cap int,
	cs ...parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, []T] {
	return parsec.Sequence[rune, Position, T](cap, cs...)
}

// Choice - searches for a combinator that works successfully on the input data.
// If one is not found, it returns an ParseError error.
//
// Alternatives are tried one by one without restoring the buffer position
// between attempts, while greedy combinators consume one item even on failure.
// Wrap each alternative in Try, otherwise a failed alternative
// eats input and breaks the following ones.
func Choice[T any](
	errMessage string,
	cs ...parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Choice(errMessage, cs...)
}

// Skip parses skip, discards its result,
// then parses body and returns the result of body.
func Skip[T any, S any](
	skip parsec.Combinator[rune, Position, S],
	body parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Skip(skip, body)
}

// SkipAfter parses body first, then parses skip and discards its result.
// It returns the result of body.
func SkipAfter[T any, S any](
	skip parsec.Combinator[rune, Position, S],
	body parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.SkipAfter(skip, body)
}

// SkipMany skips a sequence of items parsed by the skip combinator
// and then applies body. Unlike Many it allocates nothing for
// the skipped part.
func SkipMany[T any, S any](
	skip parsec.Combinator[rune, Position, S],
	body parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.SkipMany(skip, body)
}

// Padded skips a sequence of items parsed by the skip combinator
// before and after body.
func Padded[T any, S any](
	skip parsec.Combinator[rune, Position, S],
	body parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Padded(skip, body)
}
