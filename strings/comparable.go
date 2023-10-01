package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Eq - succeeds for any item which equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Eq(t rune) p.Combinator[rune, Position, rune] {
	return p.Eq[rune, Position](t)
}

// NotEq - succeeds for any item which not equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotEq(t rune) p.Combinator[rune, Position, rune] {
	return p.NotEq[rune, Position](t)
}

// OneOf - succeeds for any item which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func OneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.OneOf[rune, Position](data...)
}

// OneOfStrings - succeeds for any item which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func OneOfStrings(strs ...string) p.Combinator[rune, Position, string] {
	combs := make([]p.Combinator[rune, Position, string], len(strs))

	for i, str := range strs {
		combs[i] = Try(String(str))
	}

	return Choice(combs...)
}

// NoneOf - succeeds for any item which not included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NoneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.NoneOf[rune, Position](data...)
}

// SequenceOf - expects a sequence of elements in the buffer
// equal to the input data sequence. If expectations are not met,
// returns NothingMatched error.
func SequenceOf(data ...rune) p.Combinator[rune, Position, []rune] {
	return p.SequenceOf[rune, Position](data...)
}

// Map - Reads one element from the input buffer using the combinator,
// then uses the resulting element to obtain a value from the map cases and try to
// match it in cases map passed by first argument.
// If the value is not found then it returns NothingMatched error.
func Map[K comparable, V any](
	cases map[K]V,
	c p.Combinator[rune, Position, K],
) p.Combinator[rune, Position, V] {
	return p.Map[rune, Position, K, V](cases, c)
}

// MapStrings - Reads text from the input buffer using the combinator and
// match it in on the fly by cases map passed by first argument.
// If the value is not found then it returns NothingMatched error.
// This combinator use special trie-like structure for text matching.
func MapStrings[V any](
	cases map[string]V,
) p.Combinator[rune, Position, V] {
	tr := stringTrie(cases)

	return func(buffer p.Buffer[rune, Position]) (V, error) {
		current := tr.children
		pos := buffer.Position()

		var result *V

		for {
			r, err := buffer.Read(true)
			if err != nil {
				break
			}

			next, exists := current[r]
			if !exists {
				break
			}

			if next.end {
				result = &next.value
				pos = buffer.Position()
			}

			current = next.children
		}

		buffer.Seek(pos)

		if result == nil {
			return *new(V), p.NothingMatched
		}

		return *result, nil
	}
}

// String - read input text and match with string passed by first argument.
// If the text not matched then it returns NothingMatched error.
func String(str string) p.Combinator[rune, Position, string] {
	return func(buffer p.Buffer[rune, Position]) (string, error) {
		for _, r := range str {
			c, err := buffer.Read(true)
			if err != nil {
				return "", err
			}

			if r != c {
				return "", p.NothingMatched
			}
		}

		return str, nil
	}
}
