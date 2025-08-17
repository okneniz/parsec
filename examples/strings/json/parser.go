package json

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/common"
	"github.com/okneniz/parsec/strings"
	testHelpers "github.com/okneniz/parsec/testing"
)

var (
	notZero = strings.Range("expected digit between 1 and 9", '1', '9')
	digit   = strings.Digit("expected digit")

	whitespace = strings.Space("space")
)

func Bool() common.Combinator[rune, strings.Position, JSON] {
	return strings.MapStrings(
		"expected bool value",
		map[string]JSON{
			"true":  JSBool{true},
			"false": JSBool{false},
		},
	)
}

func Null() common.Combinator[rune, strings.Position, JSON] {
	return strings.Cast(
		strings.String("null"),
		func(_ string) (JSON, error) {
			return JSNull{}, nil
		},
	)
}

func Number_() common.Combinator[rune, strings.Position, JSON] {
	return strings.Cast(
		strings.Concat(
			1,
			strings.Count(1, "not zero", notZero),
			strings.Many(0, strings.Try(digit)),
		),
		func(x []rune) (JSON, error) {
			v, err := testHelpers.DigitsToNum(x)
			if err != nil {
				return nil, err
			}

			return JSNumber{v}, nil
		},
	)
}

func String_() common.Combinator[rune, strings.Position, JSON] {
	leftQuote := strings.Padded(
		whitespace,
		strings.Eq(`expecte double quote as start of string literal`, '"'),
	)

	rightQuote := strings.Padded(
		whitespace,
		strings.Eq(`expected double quote as end of string literal`, '"'),
	)

	parse := strings.ManyTill(
		0,
		"expected string literal",
		strings.Any(),
		rightQuote,
	)

	return func(
		buffer common.Buffer[rune, strings.Position],
	) (JSON, common.Error[strings.Position]) {
		_, err := leftQuote(buffer)
		if err != nil {
			return nil, err
		}

		s, err := parse(buffer)
		if err != nil {
			return nil, err
		}

		return JSString{string(s)}, nil
	}
}

func Value(t testing.TB) common.Combinator[rune, strings.Position, JSON] {
	var value common.Combinator[rune, strings.Position, JSON]

	keyComb := strings.Try(strings.Padded(whitespace, String_()))

	colon := strings.Padded(
		whitespace,
		strings.Eq("expected collon as separatar between object key and value", ':'),
	)

	pair := func(
		buffer common.Buffer[rune, strings.Position],
	) (*JSPair, common.Error[strings.Position]) {
		pos := buffer.Position()

		key, err := keyComb(buffer)
		if err != nil {
			return nil, err
		}

		_, err = strings.Padded(whitespace, colon)(buffer)
		if err != nil {
			return nil, err
		}

		val, err := value(buffer)
		if err != nil {
			return nil, err
		}

		ks, ok := key.(JSString)
		if !ok {
			return nil, common.NewParseError(pos, fmt.Sprintf("receive %#v as string", key))
		}

		return &JSPair{ks, val}, nil
	}

	comma := strings.Padded(
		whitespace,
		strings.Eq("expected comma as separator between key-value pairs in object", ','),
	)

	listOfPairs := strings.SepBy(0, pair, comma)

	leftBracket := strings.Eq(`left bracket as start of object"`, '{')
	rightBracket := strings.Eq(`right bracket as end of object`, '}')

	object := strings.Between(
		leftBracket,
		func(
			buffer common.Buffer[rune, strings.Position],
		) (JSON, common.Error[strings.Position]) {
			list, err := listOfPairs(buffer)
			if err != nil {
				return nil, err
			}

			m := make(map[string]JSON, len(list))
			for _, p := range list {
				m[p.key.value] = p.value
			}

			return JSObject{m}, nil
		},
		rightBracket,
	)

	leftSquare := strings.Eq(`expected left square as start of array`, '[')
	rightSquare := strings.Eq(`expected right square as end of array`, ']')

	array := strings.Between(
		leftSquare,
		func(
			buffer common.Buffer[rune, strings.Position],
		) (JSON, common.Error[strings.Position]) {
			listOfValues := strings.SepBy(0, value, comma)

			list, err := listOfValues(buffer)
			if err != nil {
				return nil, err
			}

			return JSArray{list}, nil
		},
		rightSquare,
	)

	// bool := common.Trace(t, "bool", Bool())
	// null := common.Trace(t, "null", Null())
	// num := common.Trace(t, "number", Number_())
	// str := common.Trace(t, "string", String_())
	// obj := common.Trace(t, "object", object)
	// arr := common.Trace(t, "array", array)

	bool := Bool()
	null := Null()
	num := Number_()
	str := String_()
	obj := object
	arr := array

	value = strings.Padded(
		whitespace,
		strings.Choice(
			strings.Try(bool),
			strings.Try(null),
			strings.Try(num),
			strings.Try(str),
			strings.Try(obj),
			strings.Try(arr),
		),
	)

	return value
}
