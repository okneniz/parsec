package strings

import (
	"github.com/okneniz/parsec"
)

// Concat applies the cs combinators one by one and concatenates
// their slice results into a single slice.
func Concat[O any, S any](
	cap int,
	cs ...parsec.Combinator[rune, Position, []O, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.Concat[rune, Position, O, S](cap, cs...)
}

// Sequence applies the cs combinators one by one
// and collects their results into a slice.
// It fails as soon as any of them fails.
func Sequence[O any, S any](
	cap int,
	cs ...parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.Sequence[rune, Position, O, S](cap, cs...)
}

// Choice - searches for a combinator that works successfully on the input data.
// If one is not found, it returns an ParseError error.
//
// Alternatives are tried one by one without restoring the buffer position
// between attempts, while greedy combinators consume one item even on failure.
// Wrap each alternative in Try, otherwise a failed alternative
// eats input and breaks the following ones.
func Choice[O any, S any](
	errMessage string,
	cs ...parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Choice(errMessage, cs...)
}

// Skip parses skip, discards its result,
// then parses body and returns the result of body.
func Skip[A any, O any, S any](
	skip parsec.Combinator[rune, Position, A, S],
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Skip(skip, body)
}

// SkipAfter parses body first, then parses skip and discards its result.
// It returns the result of body.
func SkipAfter[A any, O any, S any](
	skip parsec.Combinator[rune, Position, A, S],
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.SkipAfter(skip, body)
}

// SkipMany skips a sequence of items parsed by the skip combinator
// and then applies body. Unlike Many it allocates nothing for
// the skipped part.
func SkipMany[A any, O any, S any](
	skip parsec.Combinator[rune, Position, A, S],
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.SkipMany(skip, body)
}

// Padded skips a sequence of items parsed by the skip combinator
// before and after body.
func Padded[A any, O any, S any](
	skip parsec.Combinator[rune, Position, A, S],
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Padded(skip, body)
}
