package strings

import (
	"github.com/okneniz/parsec/common"
)

// Satisfy succeeds for any rune for which the supplied function f returns
// true and returns that rune. If greedy is true, the buffer keeps the
// position after reading - even when f fails, so the failed combinator
// consumes exactly one rune. Wrap it in Try if you need the position
// to be restored on failure.
func Satisfy(
	errMessage string,
	greedy bool,
	f common.Condition[rune],
) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, greedy, f)
}

// Any reads and returns the next rune from the buffer.
// It fails only at the end of the buffer.
func Any() common.Combinator[rune, Position, rune] {
	return common.Any[rune, Position]()
}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
func Try[T any](
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Try[rune, Position, T](c)
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[T any, S any, B any](
	pre common.Combinator[rune, Position, T],
	c common.Combinator[rune, Position, S],
	suf common.Combinator[rune, Position, B],
) common.Combinator[rune, Position, S] {
	return common.Between(pre, c, suf)
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF() common.Combinator[rune, Position, bool] {
	return common.EOF[rune, Position]()
}

// Cast parses the input with c and transforms the result with f.
// It fails when f returns an error.
func Cast[T any, S any](
	c common.Combinator[rune, Position, T],
	f func(T) (S, error),
) common.Combinator[rune, Position, S] {
	return common.Cast(c, f)
}

// Const consumes nothing and returns the given value.
func Const[S any](value S) common.Combinator[rune, Position, S] {
	return common.Const[rune, Position, S](value)
}

// Fail consumes nothing and always fails with the given message.
func Fail[S any](errMessage string) common.Combinator[rune, Position, S] {
	return common.Fail[rune, Position, S](errMessage)
}
