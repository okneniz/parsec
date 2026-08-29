package strings

import (
	"github.com/okneniz/parsec"
)

// Range succeeds when the next rune is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Range(
	errMessage string,
	from, to rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Range[rune, Position](errMessage, from, to)
}

// NotRange succeeds when the next rune is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotRange(
	errMessage string,
	from, to rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.NotRange[rune, Position](errMessage, from, to)
}

// Gt succeeds when the next rune is greater than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gt(
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Gt[rune, Position](errMessage, t)
}

// Gte succeeds when the next rune is greater than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gte(
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Gte[rune, Position](errMessage, t)
}

// Lt succeeds when the next rune is less than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lt(
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Lt[rune, Position](errMessage, t)
}

// Lte succeeds when the next rune is less than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lte(
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Lte[rune, Position](errMessage, t)
}
