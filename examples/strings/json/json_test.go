package json

import (
	"fmt"
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

			if !isJSEq(data, result) {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})
}

func isJSEq(x, y JSON) bool {
	if fmt.Sprintf("%T", x) != fmt.Sprintf("%T", y) {
		return false
	}

	switch xv := x.(type) {
	case JSBool:
		yv, ok := y.(JSBool)
		if !ok {
			return false
		}

		return xv.value == yv.value
	case JSString:
		yv, ok := y.(JSString)
		if !ok {
			return false
		}

		return xv.value == yv.value
	case JSNumber:
		yv, ok := y.(JSNumber)
		if !ok {
			return false
		}

		return xv.value == yv.value
	case JSNull:
		_, ok := y.(JSNull)
		return ok
	case JSArray:
		yv, ok := y.(JSArray)
		if !ok {
			return false
		}

		if len(xv.values) != len(yv.values) {
			return false
		}

		for i := 0; i < len(xv.values); i++ {
			if !isJSEq(xv.values[i], yv.values[i]) {
				return false
			}
		}

		return true
	case JSObject:
		yv, ok := y.(JSObject)
		if !ok {
			return false
		}

		if len(xv.values) != len(yv.values) {
			return false
		}

		for key, xvalue := range xv.values {
			yvalue, exists := yv.values[key]
			if !exists {
				return false
			}

			if !isJSEq(xvalue, yvalue) {
				return false
			}
		}

		return true
	}

	return false
}
