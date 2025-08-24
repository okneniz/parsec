package json

import (
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec/strings"
)

func TestJSONParsing(t *testing.T) {
	t.Parallel()

	const iterations = 1000

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	rnd := rand.New(rand.NewPCG(0, uint64(seed)))
	arb := newArbitraryJSON(rnd)
	comb := Value(t)

	t.Run("don't loose data on parsing -> serialization -> parsing chain", func(t *testing.T) {
		ohsnap.Check(t, iterations, arb, func(data JSON) bool {
			str, err := data.ToString()
			if err != nil {
				t.Log("data", data)
				t.Error(err)
				return false
			}

			result, err := strings.ParseString(str, comb)
			if err != nil {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if !JSEq(data, result) {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})
}
