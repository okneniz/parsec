package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Satisfy - succeeds for any item for which the supplied function f returns true.
// Returns the item that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy(
	greedy bool,
	f p.Condition[rune],
) p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](greedy, f)
}

// Any - returns the readed item.
func Any() p.Combinator[rune, Position, rune] {
	return p.Any[rune, Position]()
}

// Try - try to use c combinator, if it falls,
// it returns buffer to the previous position.
func Try[T any](c p.Combinator[rune, Position, T]) p.Combinator[rune, Position, T] {
	return p.Try[rune, Position, T](c)
}

// Between - parse sequence of input combinators, skip first and last results.
func Between[T any, S any, B any](
	pre p.Combinator[rune, Position, T],
	c p.Combinator[rune, Position, S],
	suf p.Combinator[rune, Position, B],
) p.Combinator[rune, Position, S] {
	return p.Between(pre, c, suf)
}

// EOF - checks that buffer reading has finished.
func EOF() p.Combinator[rune, Position, bool] {
	return p.EOF[rune, Position]()
}

// Cast - parse data by c combinator and apply to f function.
// Return result of f function.
func Cast[T any, S any](
	c p.Combinator[rune, Position, T],
	f func(T) (S, error),
) p.Combinator[rune, Position, S] {
	return p.Cast(c, f)
}

// Const - doesn't read anything, just return the input value.
func Const[S any](value S) p.Combinator[rune, Position, S] {
	return p.Const[rune, Position, S](value)
}

// Fail - doesn't read anything, just return input error.
func Fail[S any](err error) p.Combinator[rune, Position, S] {
	return p.Fail[rune, Position, S](err)
}
