package strings

import (
	"github.com/okneniz/parsec"
)

// Satisfy succeeds for any rune for which the supplied function f returns
// true and returns that rune. If greedy is true, the buffer keeps the
// position after reading - even when f fails, so the failed combinator
// consumes exactly one rune. Wrap it in Try if you need the position
// to be restored on failure.
func Satisfy(
	errMessage string,
	greedy bool,
	f parsec.Condition[rune],
) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, greedy, f)
}

// Any reads and returns the next rune from the buffer.
// It fails only at the end of the buffer.
func Any() parsec.Combinator[rune, Position, rune] {
	return parsec.Any[rune, Position]()
}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
func Try[T any](
	c parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.Try[rune, Position, T](c)
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[T any, S any, B any](
	pre parsec.Combinator[rune, Position, T],
	c parsec.Combinator[rune, Position, S],
	suf parsec.Combinator[rune, Position, B],
) parsec.Combinator[rune, Position, S] {
	return parsec.Between(pre, c, suf)
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF() parsec.Combinator[rune, Position, bool] {
	return parsec.EOF[rune, Position]()
}

// Cast parses the input with c and transforms the result with f.
// It fails when f returns an error.
func Cast[T any, S any](
	c parsec.Combinator[rune, Position, T],
	f func(T) (S, error),
) parsec.Combinator[rune, Position, S] {
	return parsec.Cast(c, f)
}

// Const consumes nothing and returns the given value.
func Const[S any](value S) parsec.Combinator[rune, Position, S] {
	return parsec.Const[rune, Position, S](value)
}

// Fail consumes nothing and always fails with the given message.
func Fail[S any](errMessage string) parsec.Combinator[rune, Position, S] {
	return parsec.Fail[rune, Position, S](errMessage)
}
