package strings

import (
	"github.com/okneniz/parsec/common"
)

// Eq - succeeds for any item which equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Eq(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Eq[rune, Position](errMessage, t)
}

// NotEq - succeeds for any item which not equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotEq(
	errMessage string,
	r rune,
) common.Combinator[rune, Position, rune] {
	return common.NotEq[rune, Position](errMessage, r)
}

// OneOf - succeeds for any item which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func OneOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, rune] {
	return common.OneOf[rune, Position](errMessage, data...)
}

// NoneOf - succeeds for any item which not included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NoneOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, rune] {
	return common.NoneOf[rune, Position](errMessage, data...)
}

// SequenceOf - expects a sequence of elements in the buffer
// equal to the input data sequence. If expectations are not met,
// returns ParseError error.
func SequenceOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, []rune] {
	return common.SequenceOf[rune, Position](errMessage, data...)
}

// Map - Reads one element from the input buffer using the combinator,
// then uses the resulting element to obtain a value from the map cases and try to
// match it in cases map passed by first argument.
// If the value is not found then it returns ParseError error.
func Map[K comparable, V any](
	errMessage string,
	cases map[K]V,
	c common.Combinator[rune, Position, K],
) common.Combinator[rune, Position, V] {
	return common.Map(errMessage, cases, c)
}

// String - read input text and match with string passed by first argument.
// If the text not matched then it returns ParseError error.
func String(errMessage, str string) common.Combinator[rune, Position, string] {
	return func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
		pos := buffer.Position()

		for _, r := range str {
			c, err := buffer.Read(true)
			if err != nil {
				return "", common.NewParseError(pos, errMessage)
			}

			if r != c {
				return "", common.NewParseError(pos, errMessage)
			}
		}

		return str, nil
	}
}

// MapStrings - Reads text from the input buffer using the combinator and
// match it in on the fly by cases map passed by first argument.
// Try to parse longest string if some of then have them have the same prefix.
// If the value is not found then it returns ParseError error.
// This combinator use special trie-like structure for text matching.
func MapStrings[V any](
	errMessage string,
	cases map[string]V,
) common.Combinator[rune, Position, V] {
	combCases := make(map[string]common.Combinator[rune, Position, V])
	for k, v := range cases {
		combCases[k] = common.Const[rune, Position, V](v)
	}

	return MapTree(errMessage, combCases)
}

// MapTree - Reads rune from the input buffer using the combinator and
// match it in on the fly by cases map passed by first argument.
// Try to parse longest prefix.
// If the value is not found then it returns ParseError error.
// This combinator use special trie-like structure for text matching.
func MapTree[T any](
	errMessage string,
	cases map[string]common.Combinator[rune, Position, T],
) common.Combinator[rune, Position, T] {
	return common.MapTree(
		errMessage,
		cases,
		func(s string) []rune { return []rune(s) },
	)
}
