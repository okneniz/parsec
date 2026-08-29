package strings

import (
	"github.com/okneniz/parsec/common"
)

// Eq succeeds when the next rune is equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Eq(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return common.Eq[rune, Position](errMessage, t)
}

// NotEq succeeds when the next rune is not equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotEq(
	errMessage string,
	r rune,
) common.Combinator[rune, Position, rune] {
	return common.NotEq[rune, Position](errMessage, r)
}

// OneOf succeeds when the next rune is one of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func OneOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, rune] {
	return common.OneOf[rune, Position](errMessage, data...)
}

// NoneOf succeeds when the next rune is none of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NoneOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, rune] {
	return common.NoneOf[rune, Position](errMessage, data...)
}

// SequenceOf expects the next runes to be equal to data in the same order
// and returns them as a slice.
func SequenceOf(
	errMessage string,
	data ...rune,
) common.Combinator[rune, Position, []rune] {
	return common.SequenceOf[rune, Position](errMessage, data...)
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[K comparable, V any](
	errMessage string,
	cases map[K]V,
	c common.Combinator[rune, Position, K],
) common.Combinator[rune, Position, V] {
	return common.Map(errMessage, cases, c)
}

// String expects the next runes to spell exactly str and returns it.
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

// MapStrings matches the input text against the cases keys on the fly,
// trying the longest possible key when some of them share a prefix,
// and returns the mapped value. It fails with errMessage when no key
// matches. Matching uses a trie-like structure, so the input is scanned
// only once.
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

// MapTree matches the input text against the cases keys on the fly,
// trying the longest possible key when some of them share a prefix,
// and applies the matched combinator to the rest of the input.
// It fails with errMessage when no key matches. Matching uses
// a trie-like structure, so the input is scanned only once.
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
