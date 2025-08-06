package strings

import (
	"github.com/okneniz/parsec/common"
)

// Satisfy - succeeds for any item for which the supplied function f returns true.
// Returns the item that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy(
	errMessage string,
	greedy bool,
	f common.Condition[rune],
) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, greedy, f)
}

// Any - returns the readed item.
func Any() common.Combinator[rune, Position, rune] {
	return common.Any[rune, Position]()
}

// Try - try to use c combinator, if it falls,
// it returns buffer to the previous position.
func Try[T any](
	c common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.Try[rune, Position, T](c)
}

// Between - parse sequence of input combinators, skip first and last results.
func Between[T any, S any, B any](
	pre common.Combinator[rune, Position, T],
	c common.Combinator[rune, Position, S],
	suf common.Combinator[rune, Position, B],
) common.Combinator[rune, Position, S] {
	return common.Between(pre, c, suf)
}

// EOF - checks that buffer reading has finished.
func EOF() common.Combinator[rune, Position, bool] {
	return common.EOF[rune, Position]()
}

// Cast - parse data by c combinator and apply to f function.
// Return result of f function.
func Cast[T any, S any](
	c common.Combinator[rune, Position, T],
	f func(T) (S, error),
) common.Combinator[rune, Position, S] {
	return common.Cast(c, f)
}

// Const - doesn't read anything, just return the input value.
func Const[S any](value S) common.Combinator[rune, Position, S] {
	return common.Const[rune, Position, S](value)
}

// Fail - doesn't read anything, just return input error.
func Fail[S any](errMessage string) common.Combinator[rune, Position, S] {
	return common.Fail[rune, Position, S](errMessage)
}
