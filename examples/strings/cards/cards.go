package cards

import (
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

// from https://www.regular-expressions.info/creditcard.html

func Visa() p.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			16,
			Count(1, Eq('4')),
			Count(12, IsDigit()),
			Optional(Count(3, IsDigit()), []rune{}),
		),
		toString,
	)
}

func Master() p.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			16,
			Choice(
				Concat(
					4,
					Count(1, Eq('5')),
					Count(1, Range('1', '5')),
					Count(2, Range('0', '9')),
				),
				Concat(
					4,
					Count(3, Eq('2')),
					Count(1, Range('1', '9')),
				),
				Concat(
					4,
					Count(2, Eq('2')),
					Count(1, Range('3', '9')),
					Count(1, Range('0', '9')),
				),
				Concat(
					4,
					Count(1, Eq('2')),
					Count(1, Range('3', '6')),
					Count(2, Range('0', '9')),
				),
				Concat(
					4,
					Count(1, Eq('2')),
					Count(1, Range('3', '6')),
					Count(2, Range('0', '9')),
				),
				Sequence(
					4,
					Eq('2'),
					Eq('7'),
					OneOf('0', '1'),
					Range('0', '9'),
				),
				SequenceOf('2', '7', '2', '0'),
			),
			Count(12, Range('0', '9')),
		),
		toString,
	)
}

func AmericanExpress() p.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			2,
			Sequence(2, Eq('3'), OneOf('4', '7')),
			Count(13, Range('0', '9')),
		),
		toString,
	)
}

func toString(x []rune) (string, error) {
	return string(x), nil
}
