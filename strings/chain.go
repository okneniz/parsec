package strings

import (
	"iter"

	"github.com/okneniz/parsec"
)

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[O any, S any](
	def O,
	c parsec.Combinator[rune, Position, O, S],
	op parsec.Combinator[rune, Position, parsec.BinaryOp[O], S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Chainl[rune, Position, O, S](def, c, op)
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[O any, S any](
	c parsec.Combinator[rune, Position, O, S],
	op parsec.Combinator[rune, Position, parsec.BinaryOp[O], S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Chainl1[rune, Position, O, S](c, op)
}

// Chainr parses zero or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// When nothing is parsed, def is returned.
func Chainr[O any, S any](
	def O,
	c parsec.Combinator[rune, Position, O, S],
	op parsec.Combinator[rune, Position, parsec.BinaryOp[O], S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Chainr[rune, Position, O, S](def, c, op)
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[O any, S any](
	c parsec.Combinator[rune, Position, O, S],
	op parsec.Combinator[rune, Position, parsec.BinaryOp[O], S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Chainr1[rune, Position, O, S](c, op)
}

// SepBy parses zero or more values of body separated by sep
// and returns them as a slice. A trailing separator is not allowed.
func SepBy[O any, A any, S any](
	cap int,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.SepBy[rune, Position, O, A, S](cap, body, sep)
}

// SepBy1 is like SepBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepBy1[O any, A any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.SepBy1[rune, Position, O, A, S](cap, errMessage, body, sep)
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[O any, A any, S any](
	cap int,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.EndBy[rune, Position, O, A, S](cap, body, sep)
}

// EndBy1 is like EndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func EndBy1[O any, A any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.EndBy1[rune, Position, O, A, S](cap, errMessage, body, sep)
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[O any, A any, S any](
	cap int,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.SepEndBy[rune, Position, O, A, S](cap, body, sep)
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[O any, A any, S any](
	cap int,
	errMessage string,
	body parsec.Combinator[rune, Position, O, S],
	sep parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.SepEndBy1[rune, Position, O, A, S](cap, errMessage, body, sep)
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[O any, A any, S any](
	cap int,
	errMessage string,
	c parsec.Combinator[rune, Position, O, S],
	end parsec.Combinator[rune, Position, A, S],
) parsec.Combinator[rune, Position, []O, S] {
	return parsec.ManyTill[rune, Position, O, A, S](cap, errMessage, c, end)
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
func Seq[O any, S any](
	c parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, iter.Seq2[O, parsec.Error[Position]], S] {
	return parsec.Seq[rune, Position, O, S](c)
}
