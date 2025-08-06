package cards

import (
	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

// from https://www.regular-expressions.info/creditcard.html

func Visa() common.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			16,
			Count(
				1,
				"expected '4'",
				Eq("expected '4'", '4'),
			),
			Count(
				12,
				"expected 12 digits",
				Digit("expected digit"),
			),
			Optional(
				Count(
					3,
					"expected 3 digits",
					Digit("expected digit"),
				),
				[]rune{},
			),
		),
		toString,
	)
}

func Master() common.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			16,
			Choice(
				Concat(
					4,
					Count(
						1,
						"expected '5'",
						Eq("expected '5'", '5'),
					),
					Count(
						1,
						"expected 'digit between '1' and '5'",
						Range("expected digit between '1' and '5'", '1', '5'),
					),
					Count(
						2,
						"expected two digits between '0' and '9'",
						Range("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						3,
						"expected '2' repeated three times",
						Eq("expected '2'", '2'),
					),
					Count(
						1,
						"expected digit between '1' and '9'",
						Range("", '1', '9'),
					),
				),
				Concat(
					4,
					Count(
						2,
						"expected '2' repeated two times",
						Eq("expected '2'", '2')),
					Count(
						1,
						"expected digit between '3' and '9'",
						Range("expected digit between '3' and '9'", '3', '9')),
					Count(
						1,
						"expected digit between '0' and '9'",
						Range("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						1,
						"expected '2'",
						Eq("expectd digit '2'", '2'),
					),
					Count(
						1,
						"expected digit between '3' and '6'",
						Range("expected digit between '3' and '6'", '3', '6'),
					),
					Count(
						2,
						"expected digits between '0' and '9' repeated two times",
						Range("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Concat(
					4,
					Count(
						1,
						"expected '2'",
						Eq("expected '2'", '2'),
					),
					Count(
						1,
						"expected digit between '3' and '6'",
						Range("expected digit betwee '3' and '6'", '3', '6'),
					),
					Count(
						2,
						"expected two digits between '0' and '9'",
						Range("expected digit between '0' and '9'", '0', '9'),
					),
				),
				Sequence(
					4,
					Eq("expected '2'", '2'),
					Eq("expected '7'", '7'),
					OneOf("expected '0' or '1'", '0', '1'),
					Range("expected digit between '0' and '9'", '0', '9'),
				),
				Sequence(
					4,
					Eq("expected '2'", '2'),
					Eq("expected '7'", '7'),
					Eq("expected '2'", '2'),
					Eq("expected '0'", '0'),
				),
			),
			Count(
				12,
				"expected digits twelve between '0' and '9'",
				Range("expected digit between '0' and '9'", '0', '9'),
			),
		),
		toString,
	)
}

func AmericanExpress() common.Combinator[rune, Position, string] {
	return Cast(
		Concat(
			2,
			Sequence(
				2,
				Eq("expected '3'", '3'),
				OneOf("expected '4' or '7'", '4', '7'),
			),
			Count(
				13,
				"expected thirteen digits between '0' and '9'",
				Range("expected digit between '0' and '9'", '0', '9'),
			),
		),
		toString,
	)
}

func toString(x []rune) (string, error) {
	return string(x), nil
}
