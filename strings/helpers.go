package strings

import (
	"github.com/okneniz/parsec/common"
)

// Parens parses body between the '(' and ')' characters
// and returns the result of body.
func Parens[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '('", '('),
		body,
		Eq("expected ')'", ')'),
	)
}

// Braces parses body between the '{' and '}' characters
// and returns the result of body.
func Braces[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '{'", '{'),
		body,
		Eq("expected '}'", '}'),
	)
}

// Angles parses body between the '<' and '>' characters
// and returns the result of body.
func Angles[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '<'", '<'),
		body,
		Eq("expected '>'", '>'),
	)
}

// Squares parses body between the '[' and ']' characters
// and returns the result of body.
func Squares[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '['", '['),
		body,
		Eq("expected ']'", ']'),
	)
}

// Semi parses the ';' character.
func Semi() common.Combinator[rune, Position, rune] {
	return Eq("expected ';'", ';')
}

// Comma parses the ',' character.
func Comma() common.Combinator[rune, Position, rune] {
	return Eq("expected ','", ',')
}

// Colon parses the ':' character.
func Colon() common.Combinator[rune, Position, rune] {
	return Eq("expected ':'", ':')
}

// Dot parses the '.' character.
func Dot() common.Combinator[rune, Position, rune] {
	return Eq("expected '.'", '.')
}

// Unsigned parses an unsigned decimal integer of type T
// (for example int or uint16). Leading zeros are accepted.
func Unsigned[T common.Integer]() common.Combinator[rune, Position, T] {
	digit := Try(Digit("digit"))
	zero := rune('0')

	return func(buffer common.Buffer[rune, Position]) (T, common.Error[Position]) {
		var result T

		token, err := digit(buffer)
		if err != nil {
			return result, err
		}

		result = T(token - zero)
		for {
			token, err = digit(buffer)
			if err != nil {
				break
			}

			result = result * 10
			result += T(token - zero)
		}

		return result, nil
	}
}

// UnsignedN parses an unsigned decimal integer of type T
// consisting of exactly n digits, for example a year or a month field
// of a fixed-width timestamp.
func UnsignedN[T common.Integer](n int, errMessage string) common.Combinator[rune, Position, T] {
	digit := Try(Digit("digit"))
	zero := rune('0')

	return func(buffer common.Buffer[rune, Position]) (T, common.Error[Position]) {
		var result T

		token, err := digit(buffer)
		if err != nil {
			return result, err
		}

		result = T(token - zero)

		for i := 0; i < n-1; i++ {
			token, err = digit(buffer)
			if err != nil {
				return 0, err
			}

			result = result * 10
			result += T(token - zero)
		}

		return result, nil
	}
}
