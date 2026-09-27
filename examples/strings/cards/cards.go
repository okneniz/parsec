package cards

import (
	"github.com/okneniz/parsec"
	. "github.com/okneniz/parsec/strings"
)

// from https://www.regular-expressions.info/creditcard.html

func Visa() parsec.Combinator[rune, Position, string, parsec.Stateless] {
	return Cast(
		Concat(
			16,
			Count(
				1,
				"expected '4'",
				Eq[parsec.Stateless]("expected '4'", '4'),
			),
			Count(
				12,
				"expected 12 digits",
				Digit[parsec.Stateless]("expected digit"),
			),
			Optional(
				Count(
					3,
					"expected 3 digits",
					Digit[parsec.Stateless]("expected digit"),
				),
				[]rune{},
			),
		),
		toString,
	)
}

func Master() parsec.Combinator[rune, Position, string, parsec.Stateless] {
	return Cast(
		Concat(
			16,
			Choice(
				"expected master card",
				Concat(
					4,
					Count(
						1,
						"expected '5'",
						Eq[parsec.Stateless]("expected '5'", '5'),
					),
					Count(
						1,
						"expected 'digit between '1' and '5'",
						Range[parsec.Stateless]("expected digit between '1' and '5'", '1', '5'),
					),
					Count(
						2,
						"expected two digits between '0' and '9'",
						Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						3,
						"expected '2' repeated three times",
						Eq[parsec.Stateless]("expected '2'", '2'),
					),
					Count(
						1,
						"expected digit between '1' and '9'",
						Range[parsec.Stateless]("", '1', '9'),
					),
				),
				Concat(
					4,
					Count(
						2,
						"expected '2' repeated two times",
						Eq[parsec.Stateless]("expected '2'", '2')),
					Count(
						1,
						"expected digit between '3' and '9'",
						Range[parsec.Stateless]("expected digit between '3' and '9'", '3', '9')),
					Count(
						1,
						"expected digit between '0' and '9'",
						Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						1,
						"expected '2'",
						Eq[parsec.Stateless]("expectd digit '2'", '2'),
					),
					Count(
						1,
						"expected digit between '3' and '6'",
						Range[parsec.Stateless]("expected digit between '3' and '6'", '3', '6'),
					),
					Count(
						2,
						"expected digits between '0' and '9' repeated two times",
						Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						1,
						"expected '2'",
						Eq[parsec.Stateless]("expected '2'", '2'),
					),
					Count(
						1,
						"expected digit between '3' and '6'",
						Range[parsec.Stateless]("expected digit betwee '3' and '6'", '3', '6'),
					),
					Count(
						2,
						"expected two digits between '0' and '9'",
						Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Sequence(
					4,
					Eq[parsec.Stateless]("expected '2'", '2'),
					Eq[parsec.Stateless]("expected '7'", '7'),
					OneOf[parsec.Stateless]("expected '0' or '1'", '0', '1'),
					Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
				),
				Sequence(
					4,
					Eq[parsec.Stateless]("expected '2'", '2'),
					Eq[parsec.Stateless]("expected '7'", '7'),
					Eq[parsec.Stateless]("expected '2'", '2'),
					Eq[parsec.Stateless]("expected '0'", '0'),
				),
			),
			Count(
				12,
				"expected digits twelve between '0' and '9'",
				Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
			),
		),
		toString,
	)
}

func AmericanExpress() parsec.Combinator[rune, Position, string, parsec.Stateless] {
	return Cast(
		Concat(
			2,
			Sequence(
				2,
				Eq[parsec.Stateless]("expected '3'", '3'),
				OneOf[parsec.Stateless]("expected '4' or '7'", '4', '7'),
			),
			Count(
				13,
				"expected thirteen digits between '0' and '9'",
				Range[parsec.Stateless]("expected digit between '0' and '9'", '0', '9'),
			),
		),
		toString,
	)
}

func toString(x []rune) (string, error) {
	return string(x), nil
}
