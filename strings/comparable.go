package strings

import (
	"github.com/okneniz/parsec"
)

// Eq succeeds when the next rune is equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func Eq[S any](
	errMessage string,
	t rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Eq[rune, Position, S](errMessage, t)
}

// NotEq succeeds when the next rune is not equal to t and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NotEq[S any](
	errMessage string,
	r rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.NotEq[rune, Position, S](errMessage, r)
}

// OneOf succeeds when the next rune is one of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func OneOf[S any](
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.OneOf[rune, Position, S](errMessage, data...)
}

// NoneOf succeeds when the next rune is none of data and returns it.
// Greedy: consumes the rune even on failure (see Try).
func NoneOf[S any](
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.NoneOf[rune, Position, S](errMessage, data...)
}

// SequenceOf expects the next runes to be equal to data in the same order
// and returns them as a slice.
func SequenceOf[S any](
	errMessage string,
	data ...rune,
) parsec.Combinator[rune, Position, []rune, S] {
	return parsec.SequenceOf[rune, Position, S](errMessage, data...)
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[K comparable, O any, S any](
	errMessage string,
	cases map[K]O,
	c parsec.Combinator[rune, Position, K, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.Map(errMessage, cases, c)
}

// String expects the next runes to spell exactly str and returns it.
func String[S any](errMessage, str string) parsec.Combinator[rune, Position, string, S] {
	return func(_ S, buffer parsec.Buffer[rune, Position]) (string, parsec.Error[Position]) {
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
func MapStrings[O any, S any](
	errMessage string,
	cases map[string]O,
) parsec.Combinator[rune, Position, O, S] {
	combCases := make(map[string]parsec.Combinator[rune, Position, O, S])
	for k, v := range cases {
		combCases[k] = parsec.Const[rune, Position, O, S](v)
	}

	return MapTree(errMessage, combCases)
}

// MapTree matches the input text against the cases keys on the fly,
// trying the longest possible key when some of them share a prefix,
// and applies the matched combinator to the rest of the input.
// It fails with errMessage when no key matches. Matching uses
// a trie-like structure, so the input is scanned only once.
func MapTree[O any, S any](
	errMessage string,
	cases map[string]parsec.Combinator[rune, Position, O, S],
) parsec.Combinator[rune, Position, O, S] {
	return parsec.MapTree(
		errMessage,
		cases,
		func(s string) []rune { return []rune(s) },
	)
}
