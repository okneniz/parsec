package timestamp

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/okneniz/oh-snap"

	"github.com/okneniz/parsec/strings"
)

func TestTimstapsParsing(t *testing.T) {
	const iterations = 100000

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	rnd := rand.New(rand.NewPCG(0, uint64(seed)))

	arbTime := ohsnap.ArbitraryTime(
		rnd,
		time.Date(1990, 5, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC),
	)

	arbLayout := ohsnap.OneOfValue(
		rnd,
		time.ANSIC,
		time.UnixDate,
		time.RFC1123,
	)

	arb := ohsnap.Combine(arbTime, arbLayout)

	parser := strings.Choice(
		strings.Try(ansic()),
		strings.Try(unixDate()),
		strings.Try(rfc1123()),
	)

	ohsnap.Check(t, iterations, arb, func(p ohsnap.Pair[time.Time, string]) bool {
		str := p.First.UTC().Format(p.Second)

		result, err := strings.ParseString(str, parser)
		if err != nil {
			t.Logf("time: %v", p.First)
			t.Logf("layout: %v", p.Second)
			t.Logf("string: %v", str)
			t.Logf("result: %v", result)
			t.Error(err)
			return false
		}

		return p.First.Truncate(time.Second).Equal(*result)
	})
}
