package json

import (
	"fmt"
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
	. "github.com/okneniz/parsec/testing"
)

var (
	notZero    = Range('1', '9')
	digit      = IsDigit()
	quote      = Padded(whitespace, Eq('"'))
	colon      = Padded(whitespace, Eq(':'))
	comma      = Padded(whitespace, Eq(','))
	whitespace = IsSpace()
)

func Bool() p.Combinator[rune, Position, JSON] {
	return Cast(
		Choice(
			Try(SequenceOf('t', 'r', 'u', 'e')),
			Try(SequenceOf('f', 'a', 'l', 's', 'e')),
		),
		func(x []rune) (JSON, error) {
			lit := string(x)
			return &JSBool{lit == "true"}, nil
		},
	)
}

func Null() p.Combinator[rune, Position, JSON] {
	return Cast(
		SequenceOf('n', 'u', 'l', 'l'),
		func(_ []rune) (JSON, error) {
			return new(JSNull), nil
		},
	)
}

func Number() p.Combinator[rune, Position, JSON] {
	return Cast(
		Concat(
			1,
			Count(1, notZero),
			Many(0, Try(digit)),
		),
		func(x []rune) (JSON, error) {
			v, err := DigitsToNum(x)
			if err != nil {
				return nil, err
			}

			return &JSNumber{v}, nil
		},
	)
}

func String_() p.Combinator[rune, Position, JSON] {
	return func(buffer p.Buffer[rune, Position]) (JSON, error) {
		_, err := quote(buffer)
		if err != nil {
			return nil, err
		}

		s, err := ManyTill(0, Any(), quote)(buffer) // TODO : move upper
		if err != nil {
			return nil, err
		}

		return &JSString{string(s)}, nil
	}
}

func Value(t testing.TB) p.Combinator[rune, Position, JSON] {
	var value p.Combinator[rune, Position, JSON]

	keyComb := Try(Padded(whitespace, String_()))

	pair := func(buffer p.Buffer[rune, Position]) (*JSPair, error) {
		key, err := keyComb(buffer)
		if err != nil {
			return nil, err
		}

		_, err = Padded(whitespace, colon)(buffer)
		if err != nil {
			return nil, err
		}

		val, err := value(buffer)
		if err != nil {
			return nil, err
		}

		ks, ok := key.(*JSString)
		if !ok {
			return nil, fmt.Errorf("receive %#v as string", key)
		}

		return &JSPair{*ks, val}, nil
	}

	listOfPairs := SepBy(0, pair, comma)

	object := Between(
		Eq('{'),
		func(buffer p.Buffer[rune, Position]) (JSON, error) {

			list, err := listOfPairs(buffer)
			if err != nil {
				return nil, err
			}

			m := make(map[string]JSON, len(list))
			for _, p := range list {
				m[p.key.value] = p.value
			}

			return &JSObject{m}, nil
		},
		Eq('}'),
	)

	array := Between(
		Eq('['),
		func(buffer p.Buffer[rune, Position]) (JSON, error) {
			listOfValues := SepBy(0, value, comma)

			list, err := listOfValues(buffer)
			if err != nil {
				return nil, err
			}

			return &JSArray{list}, nil
		},
		Eq(']'),
	)

	bool := Trace(t, "bool", Bool())
	null := Trace(t, "null", Null())
	num := Trace(t, "number", Number())
	str := Trace(t, "string", String_())
	obj := Trace(t, "object", object)
	arr := Trace(t, "array", array)

	value = Padded(
		whitespace,
		Choice(
			Try(bool),
			Try(null),
			Try(num),
			Try(str),
			Try(obj),
			Try(arr), // TODO : remove last try?
		),
	)

	return value
}
