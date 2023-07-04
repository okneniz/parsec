package cards

import (
	. "git.sr.ht/~okneniz/parsec/strings"
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
	"time"
)

func TestCards(t *testing.T) {
	digit := Range('0', '9')
	notDigit := NotRange('0', '9')

	// from https://www.regular-expressions.info/creditcard.html

	visa := Trace(t, "visa",
		Cast(
			Sequence(
				16,
				Count(1, Eq('4')),
				Count(12, digit),
				Optional(Count(3, digit), []rune{}),
			),
			toString,
		),
	)

	master := Trace(t, "master",
		Cast(
			Sequence(
				15,
				Choice(
					Trace(t, "master 1",
						Concat(
							4,
							Count(1, Eq('5')),
							Count(1, Range('1', '5')),
							Count(2, Range('0', '9')),
						),
					),
					Trace(t, "master 2",
						Concat(
							4,
							Count(3, Eq('2')),
							Count(1, Range('1', '9')),
						),
					),
					Trace(t, "master 3",
						Concat(
							4,
							Count(2, Eq('2')),
							Count(1, Range('3', '9')),
							Count(1, Range('0', '9')),
						),
					),
					Trace(t, "master 4",
						Concat(
							4,
							Count(1, Eq('2')),
							Count(1, Range('3', '6')),
							Count(2, Range('0', '9')),
						),
					),
					Trace(t, "master 5",
						Concat(
							4,
							Count(1, Eq('2')),
							Count(1, Range('3', '6')),
							Count(2, Range('0', '9')),
						),
					),
					Trace(t, "master 5",
						Sequence(
							4,
							Eq('2'),
							Eq('7'),
							OneOf('0', '1'),
							Range('0', '9'),
						),
					),
					Trace(t, "master 6",
						SequenceOf('2', '7', '2', '0'),
					),
				),
				Count(12, Range('0', '9')),
			),
			toString,
		),
	)

	americanExpress := Trace(t, "american express",
		Cast(
			Sequence(
				2,
				Sequence(
					2,
					Eq('3'),
					OneOf('4', '7'),
				),
				Count(13, Range('0', '9')),
			),
			toString,
		),
	)

	cards := Choice(
		Try(visa),
		Try(master),
		Try(americanExpress),
	)

	noice := Many(10, Try(notDigit))
	comb := Skip(noice, SepEndBy(4, cards, noice))

	cardNums := []string{
		"4111111111111111",
		"4012888888881881",
		"4222222222222",
		"5555555555554444",
		"5105105105105100",
		"378282246310005",
		"371449635398431",
		"378734493671000",
	}

	seed := time.Now().UnixNano()
	t.Log("seed: ", seed)
	shuffle := Shuffler[string](seed)

	input := JoinBy(
		Noicer(seed, '0', '9'),
		shuffle(cardNums)...,
	)

	t.Log("input:")
	t.Logf("%#v", input)

	result, err := ParseString(input, comb)
	Check(t, err)
	AssertSlice(t, Sorted(result...), Sorted(cardNums...))
}

func toString(xs [][]rune) (string, error) {
	s := ""
	for _, x := range xs {
		s += string(x)
	}
	return s, nil
}
