package strings

import (
	"github.com/okneniz/parsec"
)

// Parens parses body between the '(' and ')' characters
// and returns the result of body.
func Parens[O any, S any](
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return Between(
		Eq[S]("expected '('", '('),
		body,
		Eq[S]("expected ')'", ')'),
	)
}

// Braces parses body between the '{' and '}' characters
// and returns the result of body.
func Braces[O any, S any](
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return Between(
		Eq[S]("expected '{'", '{'),
		body,
		Eq[S]("expected '}'", '}'),
	)
}

// Angles parses body between the '<' and '>' characters
// and returns the result of body.
func Angles[O any, S any](
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return Between(
		Eq[S]("expected '<'", '<'),
		body,
		Eq[S]("expected '>'", '>'),
	)
}

// Squares parses body between the '[' and ']' characters
// and returns the result of body.
func Squares[O any, S any](
	body parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return Between(
		Eq[S]("expected '['", '['),
		body,
		Eq[S]("expected ']'", ']'),
	)
}

// Semi parses the ';' character.
func Semi[S any]() parsec.Combinator[rune, Position, rune, S] {
	return Eq[S]("expected ';'", ';')
}

// Comma parses the ',' character.
func Comma[S any]() parsec.Combinator[rune, Position, rune, S] {
	return Eq[S]("expected ','", ',')
}

// Colon parses the ':' character.
func Colon[S any]() parsec.Combinator[rune, Position, rune, S] {
	return Eq[S]("expected ':'", ':')
}

// Dot parses the '.' character.
func Dot[S any]() parsec.Combinator[rune, Position, rune, S] {
	return Eq[S]("expected '.'", '.')
}

// Unsigned parses an unsigned decimal integer of type O
// (for example int or uint16). Leading zeros are accepted.
func Unsigned[O parsec.Integer, S parsec.Recoverable[S]]() parsec.Combinator[rune, Position, O, S] {
	digit := Try[rune, S](Digit[S]("digit"))
	zero := rune('0')

	return func(state S, buffer parsec.Buffer[rune, Position]) (O, parsec.Error[Position]) {
		var result O

		token, err := digit(state, buffer)
		if err != nil {
			return result, err
		}

		result = O(token - zero)
		for {
			token, err = digit(state, buffer)
			if err != nil {
				break
			}

			result = result * 10
			result += O(token - zero)
		}

		return result, nil
	}
}

// UnsignedN parses an unsigned decimal integer of type O
// consisting of exactly n digits, for example a year or a month field
// of a fixed-width timestamp.
func UnsignedN[O parsec.Integer, S parsec.Recoverable[S]](n int, errMessage string) parsec.Combinator[rune, Position, O, S] {
	digit := Try[rune, S](Digit[S]("digit"))
	zero := rune('0')

	return func(state S, buffer parsec.Buffer[rune, Position]) (O, parsec.Error[Position]) {
		var result O

		token, err := digit(state, buffer)
		if err != nil {
			return result, err
		}

		result = O(token - zero)

		for i := 0; i < n-1; i++ {
			token, err = digit(state, buffer)
			if err != nil {
				return 0, err
			}

			result = result * 10
			result += O(token - zero)
		}

		return result, nil
	}
}
