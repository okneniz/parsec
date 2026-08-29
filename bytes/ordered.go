package bytes

import (
	"github.com/okneniz/parsec"
)

// Range succeeds when the next byte is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Range(
	errMessage string,
	from, to byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.Range[byte, int](errMessage, from, to)
}

// NotRange succeeds when the next byte is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NotRange(
	errMessage string,
	from, to byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.NotRange[byte, int](errMessage, from, to)
}

// Gt succeeds when the next byte is greater than t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Gt(
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.Gt[byte, int](errMessage, t)
}

// Gte succeeds when the next byte is greater than or equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Gte(
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.Gte[byte, int](errMessage, t)
}

// Lt succeeds when the next byte is less than t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Lt(
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.Lt[byte, int](errMessage, t)
}

// Lte succeeds when the next byte is less than or equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Lte(
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte] {
	return parsec.Lte[byte, int](errMessage, t)
}
