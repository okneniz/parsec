package bytes

import (
	"iter"

	"github.com/okneniz/parsec"
)

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[O any, S any](
	def O,
	c parsec.Combinator[byte, int, O, S],
	op parsec.Combinator[byte, int, parsec.BinaryOp[O], S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Chainl[byte, int, O, S](def, c, op)
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[O any, S any](
	c parsec.Combinator[byte, int, O, S],
	op parsec.Combinator[byte, int, parsec.BinaryOp[O], S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Chainl1[byte, int, O, S](c, op)
}

// Chainr parses zero or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// When nothing is parsed, def is returned.
func Chainr[O any, S any](
	def O,
	c parsec.Combinator[byte, int, O, S],
	op parsec.Combinator[byte, int, parsec.BinaryOp[O], S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Chainr[byte, int, O, S](def, c, op)
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[O any, S any](
	c parsec.Combinator[byte, int, O, S],
	op parsec.Combinator[byte, int, parsec.BinaryOp[O], S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Chainr1[byte, int, O, S](c, op)
}

// SepBy parses zero or more values of body separated by sep
// and returns them as a slice. A trailing separator is not allowed.
func SepBy[O any, A any, S parsec.Recoverable[S]](
	limit int,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.SepBy[byte, int, O, A, S](limit, body, sep)
}

// SepBy1 is like SepBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepBy1[O any, A any, S parsec.Recoverable[S]](
	limit int,
	errMessage string,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.SepBy1[byte, int, O, A, S](limit, errMessage, body, sep)
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[O any, A any, S parsec.Recoverable[S]](
	limit int,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.EndBy[byte, int, O, A, S](limit, body, sep)
}

// EndBy1 is like EndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func EndBy1[O any, A any, S parsec.Recoverable[S]](
	limit int,
	errMessage string,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.EndBy1[byte, int, O, A, S](limit, errMessage, body, sep)
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[O any, A any, S parsec.Recoverable[S]](
	limit int,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.SepEndBy[byte, int, O, A, S](limit, body, sep)
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[O any, A any, S parsec.Recoverable[S]](
	limit int,
	errMessage string,
	body parsec.Combinator[byte, int, O, S],
	sep parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.SepEndBy1[byte, int, O, A, S](limit, errMessage, body, sep)
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[O any, A any, S parsec.Recoverable[S]](
	limit int,
	errMessage string,
	c parsec.Combinator[byte, int, O, S],
	end parsec.Combinator[byte, int, A, S],
) parsec.Combinator[byte, int, []O, S] {
	return parsec.ManyTill[byte, int, O, A, S](limit, errMessage, c, end)
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
	c parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, iter.Seq2[O, parsec.Error[int]], S] {
	return parsec.Seq[byte, int, O, S](c)
}
