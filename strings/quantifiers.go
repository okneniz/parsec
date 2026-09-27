package strings

import (
	"github.com/okneniz/parsec"
)

// Optional applies c and returns its result,
// or the def value when c fails. It never fails.
func Optional[O any, S any](
	c parsec.Combinator[rune, Position, O, S],
	def O,
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Optional[rune, Position, O, S](c, def)
}

// Many applies c as many times as possible and collects the results.
// It stops at the first error or at the end of the buffer and returns
// everything collected so far, possibly an empty slice.
// Wrap c in Try to stop without consuming input.
func Many[O any, S any](
	cap int,
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.Many[rune, Position, O, S](cap, c)
}

// Some is like Many but requires at least one item:
// it fails with errMessage when nothing could be parsed.
func Some[O any, S any](
	cap int,
	errMessage string,
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.Some[rune, Position, O, S](cap, errMessage, c)
}

// Count applies c exactly cap times and collects the results.
// It fails with errMessage as soon as any application fails.
func Count[O any, S any](
	cap int,
	errMessage string,
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.Count[rune, Position, O, S](cap, errMessage, c)
}
