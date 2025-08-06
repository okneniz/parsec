package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Chainl - read zero or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl[T any](
	c common.Combinator[byte, int, T],
	op common.Combinator[byte, int, func(T, T) T],
	def T,
) common.Combinator[byte, int, T] {
	return common.Chainl[byte, int, T](c, op, def)
}

// Chainl1 - read one or more occurrences of byte readed by c combinator,
// separated by data readed by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl1[T any](
	c common.Combinator[byte, int, T],
	op common.Combinator[byte, int, func(T, T) T],
) common.Combinator[byte, int, T] {
	return common.Chainl1[byte, int, T](c, op)
}

// Chainr - read zero or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr[T any](
	c common.Combinator[byte, int, T],
	op common.Combinator[byte, int, func(T, T) T],
	def T,
) common.Combinator[byte, int, T] {
	return common.Chainr[byte, int, T](c, op, def)
}

// Chainr - read one or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr1[T any](
	c common.Combinator[byte, int, T],
	op common.Combinator[byte, int, func(T, T) T],
) common.Combinator[byte, int, T] {
	return common.Chainr1[byte, int, T](c, op)
}

// SepBy - read zero or more occurrences of byte readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy[T any, S any](
	cap int,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.SepBy[byte, int, T](cap, body, sep)
}

// SepBy1 - read one or more occurrences of byte readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.SepBy1[byte, int, T](cap, errMessage, body, sep)
}

// EndBy - read zero or more occurrences of byte readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by p.
func EndBy[T any, S any](
	cap int,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.EndBy[byte, int, T](cap, body, sep)
}

// EndBy1 - read one or more occurrences of byte readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by c combinator.
func EndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.EndBy1[byte, int, T](cap, errMessage, body, sep)
}

// SepEndBy - read zero or more occurrences of byte readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy[T any, S any](
	cap int,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.SepEndBy[byte, int, T](cap, body, sep)
}

// SepEndBy1 - read one or more occurrences of byte readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy1[T any, S any](
	cap int,
	errMessage string,
	body common.Combinator[byte, int, T],
	sep common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.SepEndBy1[byte, int, T](cap, errMessage, body, sep)
}

// ManyTill - accumulate data readed by c combinator until combinantor end succeeds.
// Returns a slice of values returned by body combinator.
func ManyTill[T any, S any](
	cap int,
	c common.Combinator[byte, int, T],
	end common.Combinator[byte, int, S],
) common.Combinator[byte, int, []T] {
	return common.ManyTill[byte, int, T](cap, c, end)
}
