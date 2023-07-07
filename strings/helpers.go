package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
	"golang.org/x/exp/constraints"
)

func Parens[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('('), body, Eq(')'))
}

func Braces[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('{'), body, Eq('}'))
}

func Angles[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('<'), body, Eq('>'))
}

func Squares[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('['), body, Eq(']'))
}

func Semi() p.Combinator[rune, Position, rune] {
	return Eq(';')
}

func Comma() p.Combinator[rune, Position, rune] {
	return Eq(',')
}

func Colon() p.Combinator[rune, Position, rune] {
	return Eq(':')
}

func Dot() p.Combinator[rune, Position, rune] {
	return Eq('.')
}

func Unsigned[T constraints.Unsigned]() p.Combinator[rune, Position, T] {
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

func UnsignedN[T constraints.Unsigned](n int) p.Combinator[rune, Position, T] {
	digit := Try(Range('0', '9'))
	zero := rune('0')

	return func(buffer p.Buffer[rune, Position]) (T, error) {
		var result T

		token, err := digit(buffer)
		if err != nil {
			return result, err
		}

		result = T(token - zero)

		for i := 0; i < n - 1; i++ {
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
