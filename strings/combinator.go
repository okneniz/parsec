package strings

import (
	"github.com/okneniz/parsec"
)

// Satisfy succeeds for any rune for which the supplied function f returns
// true and returns that rune. If greedy is true, the buffer keeps the
// position after reading - even when f fails, so the failed combinator
// consumes exactly one rune. Wrap it in Try if you need the position
// to be restored on failure.
func Satisfy[S any](
	errMessage string,
	greedy bool,
	f parsec.Condition[rune],
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, greedy, f)
}

// Any reads and returns the next rune from the buffer.
// It fails only at the end of the buffer.
func Any[S any]() parsec.Combinator[rune, Position, rune, S] {
	return parsec.Any[rune, Position, S]()
}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
// The state is not rewound: a combinator that changes it must restore
// the change itself before failing.
func Try[O any, S any](
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Try[rune, Position, O, S](c)
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[A any, B any, O any, S any](
	pre parsec.Combinator[rune, Position, A, S],
	c parsec.Combinator[rune, Position, O, S],
	suf parsec.Combinator[rune, Position, B, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Between(pre, c, suf)
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF[S any]() parsec.Combinator[rune, Position, bool, S] {
	return parsec.EOF[rune, Position, S]()
}

// Cast parses the input with c and transforms the result with f.
// It fails when f returns an error.
func Cast[A any, O any, S any](
	c parsec.Combinator[rune, Position, A, S],
	f func(A) (O, error),
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Cast(c, f)
}

// Const consumes nothing and returns the given value.
func Const[O any, S any](value O) parsec.Combinator[rune, Position, O, S] {
	return parsec.Const[rune, Position, O, S](value)
}

// Fail consumes nothing and always fails with the given message.
func Fail[O any, S any](errMessage string) parsec.Combinator[rune, Position, O, S] {
	return parsec.Fail[rune, Position, O, S](errMessage)
}
