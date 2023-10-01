package cards

import (
	"testing"
	"time"

	. "github.com/okneniz/parsec/strings"
	. "github.com/okneniz/parsec/testing"
)

func TestCards(t *testing.T) {
	visa := Trace(t, "visa", Visa())
	master := Trace(t, "master", Master())
	americanExpress := Trace(t, "american express", AmericanExpress())

	cards := Choice(
		Try(visa),
		Try(master),
		Try(americanExpress),
	)

	noice := Many(10, Try(NotRange('0', '9')))
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
	shuffle := Shuffler[string](seed)
	input := JoinBy(
		Noicer(seed, '0', '9'),
		shuffle(cardNums)...,
	)

	t.Log("seed: ", seed)
	t.Log("input:")
	t.Logf("%#v", input)

	result, err := ParseString(input, comb)
	Check(t, err)
	AssertSlice(t, Sorted(result...), Sorted(cardNums...))
}
