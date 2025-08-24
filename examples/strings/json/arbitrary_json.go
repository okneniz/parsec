package json

import (
	"math"
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"
)

type arbitraryJSON struct {
	rand *rand.Rand

	arbBool   ohsnap.Arbitrary[bool]
	arbInt    ohsnap.Arbitrary[int]
	arbString ohsnap.Arbitrary[string]

	arbComplexSize ohsnap.Arbitrary[int]
	arbType        ohsnap.Arbitrary[int]

	maxDeep int
}

var _ ohsnap.Arbitrary[JSON] = &arbitraryJSON{}

func newArbitraryJSON(
	rnd *rand.Rand,
) *arbitraryJSON {
	return &arbitraryJSON{
		rand:           rnd,
		arbBool:        ohsnap.ArbitraryBool(rnd),
		arbInt:         ohsnap.ArbitraryInt(rnd, 0, math.MaxInt),
		arbString:      ohsnap.ArbitraryString(rnd, "abcdefghijklmnopqrstuvwxyz", 3, 15),
		arbComplexSize: ohsnap.ArbitraryInt(rnd, 0, 5),
		arbType:        ohsnap.ArbitraryInt(rnd, 0, 5),

		maxDeep: 5,
	}
}

func (a *arbitraryJSON) generate(deep int) JSON {
	if deep >= a.maxDeep {
		return JSNull{}
	}

	switch a.arbType.Generate() {
	case 0:
		return JSBool{a.arbBool.Generate()}
	case 1:
		return JSNumber{a.arbInt.Generate()}
	case 2:
		return JSString{a.arbString.Generate()}
	case 3:
		size := a.arbComplexSize.Generate()
		arr := make([]JSON, size)

		for i := 0; i < size; i++ {
			arr[i] = a.generate(deep + 1)
		}

		return JSArray{arr}
	case 4:
		size := a.arbComplexSize.Generate()

		m := make(map[string]JSON, size)
		for i := 0; i < size; i++ {
			key := a.arbString.Generate()
			value := a.generate(deep)
			m[key] = value
		}

		return JSObject{m}
	case 5:
		return JSNull{}
	default:
		panic("unsupported type")
	}
}

func (a *arbitraryJSON) Generate() JSON {
	return a.generate(0)
}

func (a *arbitraryJSON) Shrink(value JSON) []JSON {
	// without shrinking
	return nil
}
