package strings

import (
	"github.com/okneniz/parsec"
)

// Eq succeeds when the next rune is equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Eq(
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Eq[rune, Position](errMessage, t)
}

// NotEq succeeds when the next rune is not equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotEq(
	errMessage string,
	r rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.NotEq[rune, Position](errMessage, r)
}

// OneOf succeeds when the next rune is one of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func OneOf(
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.OneOf[rune, Position](errMessage, data...)
}

// NoneOf succeeds when the next rune is none of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NoneOf(
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, rune] {
	return parsec.NoneOf[rune, Position](errMessage, data...)
}

// SequenceOf expects the next runes to be equal to data in the same order
// and returns them as a slice.
func SequenceOf(
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, []rune] {
	return parsec.SequenceOf[rune, Position](errMessage, data...)
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[K comparable, V any](
	errMessage string,
	cases map[K]V,
	c parsec.Combinator[rune, Position, K],
) parsec.Combinator[rune, Position, V] {
	return parsec.Map(errMessage, cases, c)
}

// String expects the next runes to spell exactly str and returns it.
func String(errMessage, str string) parsec.Combinator[rune, Position, string] {
	return func(buffer parsec.Buffer[rune, Position]) (string, parsec.Error[Position]) {
		pos := buffer.Position()

		for _, r := range str {
			c, err := buffer.Read(true)
			if err != nil {
				return "", parsec.NewParseError(pos, errMessage)
			}

			if r != c {
				return "", parsec.NewParseError(pos, errMessage)
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
) parsec.Combinator[rune, Position, V] {
	combCases := make(map[string]parsec.Combinator[rune, Position, V])
	for k, v := range cases {
		combCases[k] = parsec.Const[rune, Position, V](v)
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
	cases map[string]parsec.Combinator[rune, Position, T],
) parsec.Combinator[rune, Position, T] {
	return parsec.MapTree(
		errMessage,
		cases,
		func(s string) []rune { return []rune(s) },
	)
}
