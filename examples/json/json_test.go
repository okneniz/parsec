package json

import (
	"fmt"
	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/strings"
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
)

func TestJSON(t *testing.T) {
	notZero := Range('1', '9')
	digit := Range('0', '9')
	quote := Eq('"')
	any := Any()
	colon := Eq(':')
	comma := Eq(',')
	whitespace := OneOf(' ', '\n', '\r', '\t')

	var value p.Combinator[rune, Position, JSON]

	bool := Trace(t, "bool", Cast(Choice(
		Try(SequenceOf('t', 'r', 'u', 'e')),
		Try(SequenceOf('f', 'a', 'l', 's', 'e')),
	), func(x []rune) (JSON, error) {
		lit := string(x)
		return &JSBool{lit == "true"}, nil
	}))

	null := Trace(t, "null", Cast(
		SequenceOf('n', 'u', 'l', 'l'),
		func(_ []rune) (JSON, error) {
			return new(JSNull), nil
		}),
	)

	num := Trace(t, "number", Cast(Concat(
		1,
		Count(1, notZero),
		Many(0, Try(digit)),
	), func(x []rune) (JSON, error) {
		v, err := DigitsToNum(x)
		if err != nil {
			return nil, err
		}

		return &JSNumber{v}, nil
	}))

	str := Trace(t, "string", func(buffer p.Buffer[rune, Position]) (JSON, error) {
		_, err := quote(buffer)
		if err != nil {
			return nil, err
		}

		s, err := ManyTill(0, any, quote)(buffer)
		if err != nil {
			return nil, err
		}

		return &JSString{string(s)}, nil
	})

	keyComb := Try(Padded(whitespace, str))

	pair := Trace(t, "pair", func(buffer p.Buffer[rune, Position]) (*JSPair, error) {
		key, err := keyComb(buffer)
		if err != nil {
			return nil, err
		}

		_, err = colon(buffer) // TODO : add optional white spaces
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
	})

	obj := Trace(t, "object", Between(
		Eq('{'),
		func(buffer p.Buffer[rune, Position]) (JSON, error) {
			listOfPairs := SepBy(0, pair, comma)

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
	))

	arr := Trace(t, "array", Between(
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
	))

	value = Padded(
		whitespace,
		Choice(
			Try(bool),
			Try(null),
			Try(num),
			Try(str),
			Try(obj),
			Try(arr),
		),
	)

	comb := value

	t.Parallel()

	t.Run("numbers", func(t *testing.T) {
		result, err := ParseString("1", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNumber{1})

		result, err = ParseString("777", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseString("777 ", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseString(" 777", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseString(" 777 ", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNumber{777})
	})

	t.Run("null", func(t *testing.T) {
		result, err := ParseString("null", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseString("  null", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseString("null ", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseString("\n\r null ", comb)
		Check(t, err)
		assertJSEq(t, result, &JSNull{})
	})

	t.Run("strings", func(t *testing.T) {
		result, err := ParseString(`"something"`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{"something"})

		result, err = ParseString(`"another"`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseString(`""`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{""})

		result, err = ParseString(`  "another"`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseString(`"another"   `, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseString(`    "another"  `, comb)
		Check(t, err)
		assertJSEq(t, result, &JSString{"another"})
	})

	t.Run("array", func(t *testing.T) {
		result, err := ParseString(`[]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{})

		result, err = ParseString(`[1]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSNumber{1},
			},
		})

		result, err = ParseString(`[null]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSNull{},
			},
		})

		result, err = ParseString(`["test"]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSString{"test"},
			},
		})

		result, err = ParseString(`["test",1,null]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSString{"test"},
				&JSNumber{1},
				&JSNull{},
			},
		})

		result, err = ParseString(`["test",1 , null, [ ],[ 2,  3,"4"]]`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSString{"test"},
				&JSNumber{1},
				&JSNull{},
				&JSArray{},
				&JSArray{
					[]JSON{
						&JSNumber{2},
						&JSNumber{3},
						&JSString{"4"},
					},
				},
			},
		})
	})

	t.Run("object", func(t *testing.T) {
		result, err := ParseString(`{}`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSObject{})

		result, err = ParseString(`{"foo":1}`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSObject{
			map[string]JSON{
				"foo": &JSNumber{1},
			},
		})

		result, err = ParseString(`{ "foo" :{"bar": [ 1, null]}  }`, comb)
		Check(t, err)
		assertJSEq(t, result, &JSObject{
			map[string]JSON{
				"foo": &JSObject{
					map[string]JSON{
						"bar": &JSArray{
							[]JSON{
								&JSNumber{1},
								&JSNull{},
							},
						},
					},
				},
			},
		})
	})
}

func assertJSEq(t *testing.T, expected, actual JSON) {
	t.Helper()

	if !isJSEq(expected, actual) {
		t.Errorf("expected: %#v", expected)
		t.Errorf("actual: %#v", actual)
		t.Fatal()
	}
}

func isJSEq(x, y JSON) bool {
	if fmt.Sprintf("%T", x) != fmt.Sprintf("%T", y) {
		return false
	}

	switch xv := x.(type) {
	case *JSString:
		yv, ok := y.(*JSString)
		if !ok {
			return false
		}

		return xv.value == yv.value
	case *JSNumber:
		yv, ok := y.(*JSNumber)
		if !ok {
			return false
		}

		return xv.value == yv.value
	case *JSNull:
		_, ok := y.(*JSNull)
		return ok
	case *JSArray:
		yv, ok := y.(*JSArray)
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
	case *JSObject:
		yv, ok := y.(*JSObject)
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
