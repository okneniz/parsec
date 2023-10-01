package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Chainl - read zero or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Chainl[rune, Position, T](c, op, def)
}

// Chainl1 - read one or more occurrences of data readed by c combinator,
// separated by data readed by op combinator.
// Returns a value obtained by a left associative application of
// all functions returned by op combinator to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainl1[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
) p.Combinator[rune, Position, T] {
	return p.Chainl1[rune, Position, T](c, op)
}

// Chainr - read zero or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Chainr[rune, Position, T](c, op, def)
}

// Chainr - read one or more occurrences of data readed by c combinator,
// separated by op combinator.
// Returns a value obtained by a right associative application
// of all functions returned by op to the values returned by c combinator.
// If nothing read, the value def is returned.
func Chainr1[T any](
	c p.Combinator[rune, Position, T],
	op p.Combinator[rune, Position, func(T, T) T],
) p.Combinator[rune, Position, T] {
	return p.Chainr1[rune, Position, T](c, op)
}

// SepBy - read zero or more occurrences of data readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepBy[rune, Position, T](cap, body, sep)
}

// SepBy1 - read one or more occurrences of data readed by c combinator,
// separated by sep combinator.
// Returns a slice of values returned by p.
func SepBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepBy1[rune, Position, T](cap, body, sep)
}

// EndBy - read zero or more occurrences of data readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by p.
func EndBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.EndBy[rune, Position, T](cap, body, sep)
}

// EndBy1 - read one or more occurrences of data readed by c combinator,
// separated and ended by data readed by sep combinator.
// Returns a slice of values returned by c combinator.
func EndBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.EndBy1[rune, Position, T](cap, body, sep)
}

// SepEndBy - read zero or more occurrences of data readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepEndBy[rune, Position, T](cap, body, sep)
}

// SepEndBy1 - read one or more occurrences of data readed by body combinator,
// separated and optionally ended by data readed by sep combinator.
// Returns a slice of values returned by body combinator.
func SepEndBy1[T any, S any](
	cap int,
	body p.Combinator[rune, Position, T],
	sep p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.SepEndBy1[rune, Position, T](cap, body, sep)
}

// ManyTill - accumulate data readed by c combinator until combinantor end succeeds.
// Returns a slice of values returned by body combinator.
func ManyTill[T any, S any](
	cap int,
	c p.Combinator[rune, Position, T],
	end p.Combinator[rune, Position, S],
) p.Combinator[rune, Position, []T] {
	return p.ManyTill[rune, Position, T](cap, c, end)
}
