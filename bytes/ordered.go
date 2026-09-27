package bytes

import (
	"github.com/okneniz/parsec"
)

// Range succeeds when the next byte is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Range[S any](
	errMessage string,
	from, to byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Range[byte, int, S](errMessage, from, to)
}

// NotRange succeeds when the next byte is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NotRange[S any](
	errMessage string,
	from, to byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.NotRange[byte, int, S](errMessage, from, to)
}

// Gt succeeds when the next byte is greater than t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Gt[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Gt[byte, int, S](errMessage, t)
}

// Gte succeeds when the next byte is greater than or equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Gte[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Gte[byte, int, S](errMessage, t)
}

// Lt succeeds when the next byte is less than t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Lt[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Lt[byte, int, S](errMessage, t)
}

// Lte succeeds when the next byte is less than or equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Lte[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Lte[byte, int, S](errMessage, t)
}
