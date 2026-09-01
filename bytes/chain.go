package bytes

import (
	"iter"

	"github.com/okneniz/parsec"
)

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[T any](
	def T,
	c parsec.Combinator[byte, int, T],
	op parsec.Combinator[byte, int, parsec.BinaryOp[T]],
) parsec.Combinator[byte, int, T] {
	return parsec.Chainl[byte, int, T](def, c, op)
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[T any](
	c parsec.Combinator[byte, int, T],
	op parsec.Combinator[byte, int, parsec.BinaryOp[T]],
) parsec.Combinator[byte, int, T] {
	return parsec.Chainl1[byte, int, T](c, op)
}

// Chainr parses zero or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// When nothing is parsed, def is returned.
func Chainr[T any](
	def T,
	c parsec.Combinator[byte, int, T],
	op parsec.Combinator[byte, int, parsec.BinaryOp[T]],
) parsec.Combinator[byte, int, T] {
	return parsec.Chainr[byte, int, T](def, c, op)
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[T any](
	c parsec.Combinator[byte, int, T],
	op parsec.Combinator[byte, int, parsec.BinaryOp[T]],
) parsec.Combinator[byte, int, T] {
	return parsec.Chainr1[byte, int, T](c, op)
}

// SepBy parses zero or more values of body separated by sep
// and returns them as a slice. A trailing separator is not allowed.
func SepBy[T any, S any](
	cap int,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.SepBy[byte, int, T](cap, body, sep)
}

// SepBy1 is like SepBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepBy1[T any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.SepBy1[byte, int, T](cap, errMessage, body, sep)
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[T any, S any](
	cap int,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.EndBy[byte, int, T](cap, body, sep)
}

// EndBy1 is like EndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func EndBy1[T any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.EndBy1[byte, int, T](cap, errMessage, body, sep)
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[T any, S any](
	cap int,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.SepEndBy[byte, int, T](cap, body, sep)
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[T any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[byte, int, T],
	sep parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.SepEndBy1[byte, int, T](cap, errMessage, body, sep)
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[T any, S any](
	cap int,
	errMessage string,
	c parsec.Combinator[byte, int, T],
	end parsec.Combinator[byte, int, S],
) parsec.Combinator[byte, int, []T] {
	return parsec.ManyTill[byte, int, T](cap, errMessage, c, end)
}

// Seq applies c lazily and repeatedly: running the combinator on a
// buffer returns an iterator over the results of every step. Each
// successful step yields its value with a nil error; the first
// failing step is yielded with its error and ends the sequence, and
// the end of the buffer ends it silently. The iterator advances the
// shared buffer as the consumer ranges — an early break leaves the
// buffer right after the last yielded value, ready for the next
// combinator. Wrap c in Try to keep the buffer intact after the
// failing step. Like Many, it never ends when c succeeds without
// consuming input.
func Seq[T any](
	c parsec.Combinator[byte, int, T],
) parsec.Combinator[byte, int, iter.Seq2[T, parsec.Error[int]]] {
	return parsec.Seq[byte, int, T](c)
}
