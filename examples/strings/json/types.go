package json

import (
	"fmt"
	"strings"
)

// https://github.com/cierelabs/yaml_spirit/blob/master/doc/specs/json-ebnf.txt

type JSON interface {
	ToString() (string, error)
}

type JSString struct {
	value string
}

func (j JSString) ToString() (string, error) {
	return fmt.Sprintf("\"%s\"", j.value), nil
}

type JSNumber struct {
	value int
}

func (j JSNumber) ToString() (string, error) {
	return fmt.Sprintf("%d", j.value), nil
}

type JSObject struct {
	values map[string]JSON
}

func (j JSObject) ToString() (string, error) {
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

func (j JSArray) ToString() (string, error) {
	b := new(strings.Builder)

	b.WriteString("[")

	for index, value := range j.values {
		vs, err := value.ToString()
		if err != nil {
			return "", err
		}

		b.WriteString(vs)

		if index < len(j.values)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String(), nil
}

type JSBool struct {
	value bool
}

func (j JSBool) ToString() (string, error) {
	return fmt.Sprintf("%t", j.value), nil
}

type JSNull struct{}

func (j JSNull) ToString() (string, error) {
	return "null", nil
}

type JSPair struct {
	key   JSString
	value JSON
}

func JSEq(x, y JSON) bool {
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
			if !JSEq(xv.values[i], yv.values[i]) {
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

			if !JSEq(xvalue, yvalue) {
				return false
			}
		}

		return true
	}

	return false
}
