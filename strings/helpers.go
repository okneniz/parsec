package strings

import (
	p "github.com/okneniz/parsec/common"
	"golang.org/x/exp/constraints"
)

// Parens - parse something between parens characters - '(' and ')'.
func Parens[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('('), body, Eq(')'))
}

// Braces - parse something between braces characters - '{' and '}'.
func Braces[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('{'), body, Eq('}'))
}

// Angles - parse something between angels characters - '<' and '>'.
func Angles[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('<'), body, Eq('>'))
}

// Squares - parse something between squares characters - '[' and ']'.
func Squares[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('['), body, Eq(']'))
}

// Semi - parse semi character.
func Semi() p.Combinator[rune, Position, rune] {
	return Eq(';')
}

// Comma - parse comma character.
func Comma() p.Combinator[rune, Position, rune] {
	return Eq(',')
}

// Colon - parse colon character.
func Colon() p.Combinator[rune, Position, rune] {
	return Eq(':')
}

// Dot - parse dot character.
func Dot() p.Combinator[rune, Position, rune] {
	return Eq('.')
}

// Unsigned - parse unsigned integer.
func Unsigned[T constraints.Integer]() p.Combinator[rune, Position, T] {
	digit := Try(Range('0', '9'))
	zero := rune('0')

	return func(buffer p.Buffer[rune, Position]) (T, error) {
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
func UnsignedN[T constraints.Integer](n int) p.Combinator[rune, Position, T] {
	digit := Try(Range('0', '9'))
	zero := rune('0')

	return func(buffer p.Buffer[rune, Position]) (T, error) {
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
