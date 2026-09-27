package strings

import (
	"github.com/okneniz/parsec"
)

// Range succeeds when the next rune is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Range[S any](
	errMessage string,
	from, to rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Range[rune, Position, S](errMessage, from, to)
}

// NotRange succeeds when the next rune is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotRange[S any](
	errMessage string,
	from, to rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.NotRange[rune, Position, S](errMessage, from, to)
}

// Gt succeeds when the next rune is greater than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gt[S any](
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Gt[rune, Position, S](errMessage, t)
}

// Gte succeeds when the next rune is greater than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gte[S any](
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Gte[rune, Position, S](errMessage, t)
}

// Lt succeeds when the next rune is less than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lt[S any](
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Lt[rune, Position, S](errMessage, t)
}

// Lte succeeds when the next rune is less than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lte[S any](
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Lte[rune, Position, S](errMessage, t)
}
