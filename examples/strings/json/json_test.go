package json

import (
	"fmt"
	"testing"

	. "github.com/okneniz/parsec/strings"
	. "github.com/okneniz/parsec/testing"
)

func TestJSON(t *testing.T) {
	comb := Value(t)

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
