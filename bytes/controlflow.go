package bytes

import (
	"github.com/okneniz/parsec"
)

// Concat applies the cs combinators one by one and concatenates
// their slice results into a single slice.
func Concat[T any](
	cap int,
	cs ...parsec.Combinator[byte, int, []T],
) parsec.Combinator[byte, int, []T] {
	return parsec.Concat[byte, int, T](cap, cs...)
}

// Sequence applies the cs combinators one by one
// and collects their results into a slice.
// It fails as soon as any of them fails.
func Sequence[T any](
	cap int,
	cs ...parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, []T] {
	return parsec.Sequence[byte, int, T](cap, cs...)
}

// Choice searches for a combinator that succeeds on the input data
// and returns its result; if none is found, it returns a ParseError.
// Alternatives are tried one by one without restoring the buffer position
// between attempts, while greedy combinators consume one item even on
// failure. Wrap each alternative in Try.
func Choice[T any](
	errMessage string,
	cs ...parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.Choice[byte, int, T](errMessage, cs...)
}

// Skip parses skip, discards its result,
// then parses body and returns the result of body.
func Skip[T any, S any](
	skip parsec.Combinator[byte, int, S],
	body parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.Skip(skip, body)
}

// SkipAfter parses body first, then parses skip and discards its result.
// It returns the result of body.
func SkipAfter[T any, S any](
	skip parsec.Combinator[byte, int, S],
	body parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.SkipAfter(skip, body)
}

// Padded skips a sequence of items parsed by the skip combinator
// before and after body.
func Padded[T any, S any](
	skip parsec.Combinator[byte, int, S],
	body parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, T] {
	return parsec.Padded(skip, body)
}
