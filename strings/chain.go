package strings

import (
	"github.com/okneniz/parsec/common"
)

// Chainl parses zero or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// When nothing is parsed, def is returned.
func Chainl[T any](
	def T,
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainl[rune, Position, T](def, c, op)
}

// Chainl1 parses one or more values with c separated by op
// and combines them with left associativity: ((a op b) op c).
// It fails when the first application of c fails.
func Chainl1[T any](
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainl1[rune, Position, T](c, op)
}

// Chainr parses zero or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// When nothing is parsed, def is returned.
func Chainr[T any](
	def T,
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainr[rune, Position, T](def, c, op)
}

// Chainr1 parses one or more values with c separated by op
// and combines them with right associativity: (a op (b op c)).
// It fails when the first application of c fails.
func Chainr1[T any](
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainr1[rune, Position, T](c, op)
}

// SepBy parses zero or more values of body separated by sep
// and returns them as a slice. A trailing separator is not allowed.
func SepBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepBy[rune, Position, T](cap, body, sep)
}

// SepBy1 is like SepBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepBy1[rune, Position, T](cap, errMessage, body, sep)
}

// EndBy parses zero or more values of body, each terminated by sep,
// like statements terminated by a semicolon.
func EndBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.EndBy[rune, Position, T](cap, body, sep)
}

// EndBy1 is like EndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func EndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.EndBy1[rune, Position, T](cap, errMessage, body, sep)
}

// SepEndBy parses zero or more values of body, separated by sep
// and optionally terminated by a final sep.
func SepEndBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepEndBy[rune, Position, T](cap, body, sep)
}

// SepEndBy1 is like SepEndBy but requires at least one value:
// it fails with errMessage when nothing could be parsed.
func SepEndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepEndBy1[rune, Position, T](cap, errMessage, body, sep)
}

// ManyTill collects the results of c until the end combinator succeeds;
// the end match itself is not included. It fails with errMessage
// when c fails before end matches.
func ManyTill[T any, S any](
	cap int,
	errMessage string,
	c common.Combinator[rune, Position, T],
	end common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.ManyTill[rune, Position, T](cap, errMessage, c, end)
}
