package parsec

import (
	"fmt"
	"strings"
	"testing"
)

type JSON interface {
	ToString() (string, error)
}

type JSString struct {
	value string
}

func (j *JSString) ToString() (string, error) {
	return fmt.Sprintf("\"%s\"", j.value), nil
}

type JSNumber struct {
	value int
}

func (j *JSNumber) ToString() (string, error) {
	return fmt.Sprintf("%d", j.value), nil
}

type JSObject struct {
	values map[string]JSON
}

func (j *JSObject) ToString() (string, error) {
	b := new(strings.Builder)

	b.WriteString("{")
	index := 0

	for key, value := range j.values {
		ks := fmt.Sprintf("\"%s\"", key)

		vs, err := value.ToString()
		if err != nil {
			return "", err
		}

		b.WriteString(ks)
		b.WriteString(":")
		b.WriteString(vs)

		index++

		if index < len(j.values) {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String(), nil
}

type JSArray struct {
	values []JSON
}

func (j *JSArray) ToString() (string, error) {
	b := new(strings.Builder)

	b.WriteString("[")

	for index, value := range j.values {
		vs, err := value.ToString()
		if err != nil {
			return "", err
		}

		b.WriteString(vs)

		if index < len(j.values) {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String(), nil
}

type JSBool struct {
	value bool
}

func (j *JSBool) ToString() (string, error) {
	return fmt.Sprintf("%t", j.value), nil
}

type JSNull struct{}

func (j *JSNull) ToString() (string, error) {
	return "null", nil
}

type JSPair struct {
	key   JSString
	value JSON
}

func TestJSON(t *testing.T) {
	notZero := Range(true, byte('1'), byte('9'))
	digit := Range(true, byte('0'), byte('9'))
	quote := Eq(byte('"'))
	any := Any[byte]()
	colon := Eq(byte(':'))
	comma := Eq(byte(','))
	whitespace := OneOf(byte(' '), byte('\n'), byte('\r'), byte('\t'))

	var value Combinator[byte, JSON]

	bool := Trace(t, "bool", Cast(Choice(
		Try(Sequence(
			4,
			Eq(byte('t')),
			Eq(byte('r')),
			Eq(byte('u')),
			Eq(byte('e')),
		)),
		Try(Sequence(
			5,
			Eq(byte('f')),
			Eq(byte('a')),
			Eq(byte('l')),
			Eq(byte('s')),
			Eq(byte('e')),
		)),
	), func(x []byte) (JSON, error) {
		lit := string(x)
		return &JSBool{lit == "true"}, nil
	}))

	null := Trace(t, "null", Cast(Sequence(
		4,
		Eq(byte('n')),
		Eq(byte('u')),
		Eq(byte('l')),
		Eq(byte('l')),
	), func(_ []byte) (JSON, error) {
		return new(JSNull), nil
	}))

	num := Trace(t, "number", Cast(Concat(
		1,
		Count(1, notZero),
		Many(0, Try(digit)),
	),
		func(x []byte) (JSON, error) {
			v, err := digitsToNum(x)
			if err != nil {
				return nil, err
			}

			return &JSNumber{v}, nil
		}))

	str := Trace(t, "string", func(buffer Buffer[byte]) (JSON, error) {
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

	pair := Trace(t, "pair", func(buffer Buffer[byte]) (*JSPair, error) {
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
		Eq(byte('{')),
		func(buffer Buffer[byte]) (JSON, error) {
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
		Eq(byte('}')),
	))

	arr := Trace(t, "array", Between(
		Eq(byte('[')),
		func(buffer Buffer[byte]) (JSON, error) {
			listOfValues := SepBy(0, value, comma)

			list, err := listOfValues(buffer)
			if err != nil {
				return nil, err
			}

			return &JSArray{list}, nil
		},
		Eq(byte(']')),
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
		result, err := ParseBytes([]byte("1"), comb)
		check(t, err)
		assertJSEq(t, result, &JSNumber{1})

		result, err = ParseBytes([]byte("777"), comb)
		check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseBytes([]byte("777 "), comb)
		check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseBytes([]byte(" 777"), comb)
		check(t, err)
		assertJSEq(t, result, &JSNumber{777})

		result, err = ParseBytes([]byte(" 777 "), comb)
		check(t, err)
		assertJSEq(t, result, &JSNumber{777})
	})

	t.Run("null", func(t *testing.T) {
		result, err := ParseBytes([]byte("null"), comb)
		check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseBytes([]byte("  null"), comb)
		check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseBytes([]byte("null "), comb)
		check(t, err)
		assertJSEq(t, result, &JSNull{})

		result, err = ParseBytes([]byte("\n\r null "), comb)
		check(t, err)
		assertJSEq(t, result, &JSNull{})
	})

	t.Run("strings", func(t *testing.T) {
		result, err := ParseBytes([]byte(`"something"`), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{"something"})

		result, err = ParseBytes([]byte(`"another"`), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseBytes([]byte(`""`), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{""})

		result, err = ParseBytes([]byte(`  "another"`), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseBytes([]byte(`"another"   `), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{"another"})

		result, err = ParseBytes([]byte(`    "another"  `), comb)
		check(t, err)
		assertJSEq(t, result, &JSString{"another"})
	})

	t.Run("array", func(t *testing.T) {
		result, err := ParseBytes([]byte(`[]`), comb)
		check(t, err)
		assertJSEq(t, result, &JSArray{})

		result, err = ParseBytes([]byte(`[1]`), comb)
		check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSNumber{1},
			},
		})

		result, err = ParseBytes([]byte(`[null]`), comb)
		check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSNull{},
			},
		})

		result, err = ParseBytes([]byte(`["test"]`), comb)
		check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSString{"test"},
			},
		})

		result, err = ParseBytes([]byte(`["test",1,null]`), comb)
		check(t, err)
		assertJSEq(t, result, &JSArray{
			[]JSON{
				&JSString{"test"},
				&JSNumber{1},
				&JSNull{},
			},
		})

		result, err = ParseBytes([]byte(`["test",1 , null, [ ],[ 2,  3,"4"]
		]`), comb)
		check(t, err)
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
		result, err := ParseBytes([]byte(`{}`), comb)
		check(t, err)
		assertJSEq(t, result, &JSObject{})

		result, err = ParseBytes([]byte(`{"foo":1}`), comb)
		check(t, err)
		assertJSEq(t, result, &JSObject{
			map[string]JSON{
				"foo": &JSNumber{1},
			},
		})

		result, err = ParseBytes([]byte(`{ "foo" :{"bar": [ 1, null]}  }`), comb)
		check(t, err)
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
