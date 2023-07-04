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
