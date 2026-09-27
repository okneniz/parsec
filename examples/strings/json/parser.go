package json

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

var (
	notZero = strings.Range[parsec.Stateless]("expected digit between 1 and 9", '1', '9')
	digit   = strings.Digit[parsec.Stateless]("expected digit")

	whitespace = strings.Space[parsec.Stateless]("space")
)

func Bool() parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless] {
	return strings.MapStrings[JSON, parsec.Stateless](
		"expected bool value",
		map[string]JSON{
			"true":  JSBool{true},
			"false": JSBool{false},
		},
	)
}

func Null() parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless] {
	return strings.Cast(
		strings.String[parsec.Stateless]("expected null value", "null"),
		func(_ string) (JSON, error) {
			return JSNull{}, nil
		},
	)
}

func Number_() parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless] {
	return strings.Cast(
		strings.Concat(
			1,
			strings.Count(1, "not zero", notZero),
			strings.Many(0, strings.Try(digit)),
		),
		func(ds []rune) (JSON, error) {
			v := 0

			for i, x := range ds {
				v += int(x - '0')

				if len(ds)-1 != i {
					v *= 10
				}
			}

			return JSNumber{v}, nil
		},
	)
}

func String_() parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless] {
	leftQuote := strings.Padded(
		whitespace,
		strings.Eq[parsec.Stateless](`expecte double quote as start of string literal`, '"'),
	)

	rightQuote := strings.Padded(
		whitespace,
		strings.Eq[parsec.Stateless](`expected double quote as end of string literal`, '"'),
	)

	parse := strings.ManyTill(
		0,
		"expected string literal",
		strings.Any[parsec.Stateless](),
		rightQuote,
	)

	return func(
		state parsec.Stateless,
		buffer parsec.Buffer[rune, strings.Position],
	) (JSON, parsec.Error[strings.Position]) {
		_, err := leftQuote(state, buffer)
		if err != nil {
			return nil, err
		}

		s, err := parse(state, buffer)
		if err != nil {
			return nil, err
		}

		return JSString{string(s)}, nil
	}
}

func Value(t testing.TB) parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless] {
	var value parsec.Combinator[rune, strings.Position, JSON, parsec.Stateless]

	keyComb := strings.Try(strings.Padded(whitespace, String_()))

	colon := strings.Padded(
		whitespace,
		strings.Eq[parsec.Stateless]("expected collon as separatar between object key and value", ':'),
	)

	pair := func(
		state parsec.Stateless,
		buffer parsec.Buffer[rune, strings.Position],
	) (*JSPair, parsec.Error[strings.Position]) {
		pos := buffer.Position()

		key, err := keyComb(state, buffer)
		if err != nil {
			return nil, err
		}

		_, err = strings.Padded(whitespace, colon)(state, buffer)
		if err != nil {
			return nil, err
		}

		val, err := value(state, buffer)
		if err != nil {
			return nil, err
		}

		ks, ok := key.(JSString)
		if !ok {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("receive %#v as string", key))
		}

		return &JSPair{ks, val}, nil
	}

	comma := strings.Padded(
		whitespace,
		strings.Eq[parsec.Stateless]("expected comma as separator between key-value pairs in object", ','),
	)

	listOfPairs := strings.SepBy(0, pair, comma)

	leftBracket := strings.Eq[parsec.Stateless](`left bracket as start of object"`, '{')
	rightBracket := strings.Eq[parsec.Stateless](`right bracket as end of object`, '}')

	object := strings.Between(
		leftBracket,
		func(
			state parsec.Stateless,
			buffer parsec.Buffer[rune, strings.Position],
		) (JSON, parsec.Error[strings.Position]) {
			list, err := listOfPairs(state, buffer)
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

	leftSquare := strings.Eq[parsec.Stateless](`expected left square as start of array`, '[')
	rightSquare := strings.Eq[parsec.Stateless](`expected right square as end of array`, ']')

	array := strings.Between(
		leftSquare,
		func(
			state parsec.Stateless,
			buffer parsec.Buffer[rune, strings.Position],
		) (JSON, parsec.Error[strings.Position]) {
			listOfValues := strings.SepBy(0, value, comma)

			list, err := listOfValues(state, buffer)
			if err != nil {
				return nil, err
			}

			return JSArray{list}, nil
		},
		rightSquare,
	)

	// bool := parsec.Trace(t, "bool", Bool())
	// null := parsec.Trace(t, "null", Null())
	// num := parsec.Trace(t, "number", Number_())
	// str := parsec.Trace(t, "string", String_())
	// obj := parsec.Trace(t, "object", object)
	// arr := parsec.Trace(t, "array", array)

	bool := Bool()
	null := Null()
	num := Number_()
	str := String_()
	obj := object
	arr := array

	value = strings.Padded(
		whitespace,
		strings.Choice(
			"expected JSON value",
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
