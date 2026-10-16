package bytes

import (
	"github.com/okneniz/parsec"
)

// Concat applies the cs combinators one by one and concatenates
// their slice results into a single slice.
func Concat[O any, S any](
	cap int,
	cs ...parsec.Combinator[byte, int, []O, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.Concat[byte, int, O, S](cap, cs...)
}

// Sequence applies the cs combinators one by one
// and collects their results into a slice.
// It fails as soon as any of them fails.
func Sequence[O any, S any](
	cap int,
	cs ...parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.Sequence[byte, int, O, S](cap, cs...)
}

// Choice searches for a combinator that succeeds on the input data
// and returns its result; if none is found, it returns a ParseError.
// Alternatives are tried one by one without restoring the buffer position
// between attempts, while greedy combinators consume one item even on
// failure. Wrap each alternative in Try.
func Choice[O any, S any](
	errMessage string,
	cs ...parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Choice[byte, int, O, S](errMessage, cs...)
}

// Skip parses skip, discards its result,
// then parses body and returns the result of body.
func Skip[A any, O any, S any](
	skip parsec.Combinator[byte, int, A, S],
	body parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Skip(skip, body)
}

// SkipAfter parses body first, then parses skip and discards its result.
// It returns the result of body.
func SkipAfter[A any, O any, S any](
	skip parsec.Combinator[byte, int, A, S],
	body parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.SkipAfter(skip, body)
}

// Padded skips a sequence of items parsed by the skip combinator
// before and after body.
func Padded[A any, O any, S parsec.Recoverable[S]](
	skip parsec.Combinator[byte, int, A, S],
	body parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Padded(skip, body)
}
