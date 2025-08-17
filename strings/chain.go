package strings

import (
	"github.com/okneniz/parsec/common"
)

// Chainl - read zero or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl[T any](
	def T,
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainl[rune, Position, T](def, c, op)
}

// Chainl1 - read one or more occurrences of data readed by c combinator,
// separated by data readed by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl1[T any](
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainl1[rune, Position, T](c, op)
}

// Chainr - read zero or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr[T any](
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
	def T,
) common.Combinator[rune, Position, T] {
	return common.Chainr[rune, Position, T](def, c, op)
}

// Chainr - read one or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr1[T any](
	c common.Combinator[rune, Position, T],
	op common.Combinator[rune, Position, common.BinaryOp[T]],
) common.Combinator[rune, Position, T] {
	return common.Chainr1[rune, Position, T](c, op)
}

// SepBy - read zero or more occurrences of data readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepBy[rune, Position, T](cap, body, sep)
}

// SepBy1 - read one or more occurrences of data readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepBy1[rune, Position, T](cap, errMessage, body, sep)
}

// EndBy - read zero or more occurrences of data readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by p.
func EndBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.EndBy[rune, Position, T](cap, body, sep)
}

// EndBy1 - read one or more occurrences of data readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by c combinator.
func EndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.EndBy1[rune, Position, T](cap, errMessage, body, sep)
}

// SepEndBy - read zero or more occurrences of data readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy[T any, S any](
	cap int,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepEndBy[rune, Position, T](cap, body, sep)
}

// SepEndBy1 - read one or more occurrences of data readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[rune, Position, T],
	sep common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.SepEndBy1[rune, Position, T](cap, errMessage, body, sep)
}

// ManyTill - accumulate data readed by c combinator until combinantor end succeeds.
// Returns a slice of values returned by body combinator.
func ManyTill[T any, S any](
	cap int,
	errMessage string,
	c common.Combinator[rune, Position, T],
	end common.Combinator[rune, Position, S],
) common.Combinator[rune, Position, []T] {
	return common.ManyTill[rune, Position, T](cap, errMessage, c, end)
}
