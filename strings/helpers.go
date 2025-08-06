package strings

import (
	"golang.org/x/exp/constraints"

	"github.com/okneniz/parsec/common"
)

// Parens - parse something between parens characters - '(' and ')'.
func Parens[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '('", '('),
		body,
		Eq("expected ')'", ')'),
	)
}

// Braces - parse something between braces characters - '{' and '}'.
func Braces[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '{'", '{'),
		body,
		Eq("expected '}'", '}'),
	)
}

// Angles - parse something between angels characters - '<' and '>'.
func Angles[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '<'", '<'),
		body,
		Eq("expected '>'", '>'),
	)
}

// Squares - parse something between squares characters - '[' and ']'.
func Squares[T any](
	body common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return Between(
		Eq("expected '['", '['),
		body,
		Eq("expected ']'", ']'),
	)
}

// Semi - parse semi character.
func Semi() common.Combinator[rune, Position, rune] {
	return Eq("expected ';'", ';')
}

// Comma - parse comma character.
func Comma() common.Combinator[rune, Position, rune] {
	return Eq("expected ','", ',')
}

// Colon - parse colon character.
func Colon() common.Combinator[rune, Position, rune] {
	return Eq("expected ':'", ':')
}

// Dot - parse dot character.
func Dot() common.Combinator[rune, Position, rune] {
	return Eq("expected '.'", '.')
}

// Unsigned - parse unsigned integer.
func Unsigned[T constraints.Integer]() common.Combinator[rune, Position, T] {
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

// UnsignedN - parse unsigned integer with N count of digits.
func UnsignedN[T constraints.Integer](n int, errMessage string) common.Combinator[rune, Position, T] {
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
