package strings

import (
	"github.com/okneniz/parsec/common"
)

// Range succeeds when the next rune is inside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Range(
	errMessage string,
	from, to rune,
) common.Combinator[rune, Position, rune] {
	return common.Range[rune, Position](errMessage, from, to)
}

// NotRange succeeds when the next rune is outside the inclusive
// range [from, to] and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotRange(
	errMessage string,
	from, to rune,
) common.Combinator[rune, Position, rune] {
	return common.NotRange[rune, Position](errMessage, from, to)
}

// Gt succeeds when the next rune is greater than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gt(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Gt[rune, Position](errMessage, t)
}

// Gte succeeds when the next rune is greater than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Gte(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Gte[rune, Position](errMessage, t)
}

// Lt succeeds when the next rune is less than t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lt(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Lt[rune, Position](errMessage, t)
}

// Lte succeeds when the next rune is less than or equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Lte(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Lte[rune, Position](errMessage, t)
}
