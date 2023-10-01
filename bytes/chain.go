package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Chainl - read zero or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Chainl[byte, int, T](c, op, def)
}

// Chainl1 - read one or more occurrences of byte readed by c combinator,
// separated by data readed by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl1[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
) p.Combinator[byte, int, T] {
	return p.Chainl1[byte, int, T](c, op)
}

// Chainr - read zero or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Chainr[byte, int, T](c, op, def)
}

// Chainr - read one or more occurrences of byte readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr1[T any](
	c p.Combinator[byte, int, T],
	op p.Combinator[byte, int, func(T, T) T],
) p.Combinator[byte, int, T] {
	return p.Chainr1[byte, int, T](c, op)
}

// SepBy - read zero or more occurrences of byte readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepBy[byte, int, T](cap, body, sep)
}

// SepBy1 - read one or more occurrences of byte readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepBy1[byte, int, T](cap, body, sep)
}

// EndBy - read zero or more occurrences of byte readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by p.
func EndBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.EndBy[byte, int, T](cap, body, sep)
}

// EndBy1 - read one or more occurrences of byte readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by c combinator.
func EndBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.EndBy1[byte, int, T](cap, body, sep)
}

// SepEndBy - read zero or more occurrences of byte readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepEndBy[byte, int, T](cap, body, sep)
}

// SepEndBy1 - read one or more occurrences of byte readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy1[T any, S any](
	cap int,
	body p.Combinator[byte, int, T],
	sep p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.SepEndBy1[byte, int, T](cap, body, sep)
}

// ManyTill - accumulate data readed by c combinator until combinantor end succeeds.
// Returns a slice of values returned by body combinator.
func ManyTill[T any, S any](
	cap int,
	c p.Combinator[byte, int, T],
	end p.Combinator[byte, int, S],
) p.Combinator[byte, int, []T] {
	return p.ManyTill[byte, int, T](cap, c, end)
}
