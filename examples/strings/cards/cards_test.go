package cards

import (
	"maps"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	. "github.com/okneniz/parsec/strings"
)

func TestCards(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	rnd := rand.New(rand.NewPCG(0, uint64(seed)))

	cards := Choice(
		"expected card number",
		Try(Master()),
		Try(Visa()),
		Try(AmericanExpress()),
	)

	noice := Many(10, Try(NotRange("expected not digit", '0', '9')))
	parse := Skip(noice, SepEndBy(0, cards, noice))

	arbNoice := ohsnap.ArbitrarySlice(
		rnd,
		ohsnap.ArbitraryString(rnd, "abcdefghijklmnopqrstuvwxyz", 1, 100),
		10,
		20,
	)

	cardsNums := map[string]struct{}{
		"4111111111111111": {},
		"4012888888881881": {},
		"5555555555554444": {},
		"5105105105105100": {},
		"378282246310005":  {},
		"371449635398431":  {},
		"378734493671000":  {},
	}

	const expectedSize = 4

	arbCard := ohsnap.OneOfValue(rnd, slices.Collect(maps.Keys(cardsNums))...)
	arbCards := ohsnap.ArbitrarySlice(rnd, arbCard, expectedSize, expectedSize)

	arb := ohsnap.Map(
		ohsnap.Combine(arbNoice, arbCards),
		func(x ohsnap.Pair[[]string, []string]) string {
			data := append(x.First, x.Second...)

			rnd.Shuffle(len(data), func(i, j int) {
				data[i], data[j] = data[j], data[i]
			})

			return strings.Join(data, "")
		},
	)

	ohsnap.Check(t, 100_000, arb, func(input string) bool {
		output, err := ParseString(input, parse)
		if err != nil {
			t.Log("input:", input)
			t.Log("output:", output)
			t.Error(err)
			return false
		}

		if len(output) != expectedSize {
			t.Log("input:", input)
			t.Log("output:", output)
			t.Error("unexpected len of result", len(output), expectedSize)
			return false
		}

		for _, x := range output {
			if _, exists := cardsNums[x]; !exists {
				t.Error("invalid card number", x)
				return false
			}
		}

		return true
	})
}
