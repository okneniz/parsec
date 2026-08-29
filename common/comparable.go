package common

// Eq succeeds when the next item is equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Eq[T comparable, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return t == x
	})
}

// NotEq succeeds when the next item is not equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func NotEq[T comparable, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return t != x
	})
}

// OneOf succeeds when the next item is one of data and returns it.
// Greedy: consumes the item even on failure (see Try).
func OneOf[T comparable, P any](
	errMessage string,
	data ...T,
) Combinator[T, P, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P](errMessage, true, func(x T) bool {
		_, exists := m[x]
		return exists
	})
}

// NoneOf succeeds when the next item is none of data and returns it.
// Greedy: consumes the item even on failure (see Try).
func NoneOf[T comparable, P any](
	errMessage string,
	data ...T,
) Combinator[T, P, T] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P](errMessage, true, func(x T) bool {
		_, exists := m[x]
		return !exists
	})
}

// SequenceOf expects the next items to be equal to data in the same order
// and returns them as a slice.
func SequenceOf[T comparable, P any](
	errMessage string,
	data ...T,
) Combinator[T, P, []T] {
	return func(buffer Buffer[T, P]) ([]T, Error[P]) {
		pos := buffer.Position()

		result := make([]T, 0, len(data))

		for _, x := range data {
			token, err := buffer.Read(true)
			if err != nil {
				return nil, NewParseError(pos, errMessage)
			}

			if x != token {
				return nil, NewParseError(pos, errMessage)
			}

			result = append(result, token)
		}

		return result, nil
	}
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[T any, P any, K comparable, V any](
	errMessage string,
	cases map[K]V,
	c Combinator[T, P, K],
) Combinator[T, P, V] {
	var null V

	return func(buffer Buffer[T, P]) (V, Error[P]) {
		pos := buffer.Position()

		token, err := c(buffer)
		if err != nil {
			return null, err
		}

		result, exists := cases[token]
		if !exists {
			return null, NewParseError(pos, errMessage)
		}

		return result, nil
	}
}

// MapAs parses a key with the comb combinator, looks the key up in cases
// and applies the matched combinator to the rest of the input.
// It fails with errMessage when the key is not found.
func MapAs[T any, P any, K comparable, V any](
	errMessage string,
	cases map[K]Combinator[T, P, V],
	comb Combinator[T, P, K],
) Combinator[T, P, V] {
	var null V

	return func(buffer Buffer[T, P]) (V, Error[P]) {
		pos := buffer.Position()

		key, err := comb(buffer)
		if err != nil {
			return null, err
		}

		parseValue, exists := cases[key]
		if !exists {
			return null, NewParseError(pos, errMessage)
		}

		return parseValue(buffer)
	}
}

// MapTree matches the input against the keys of cases using a
// longest-prefix trie (keys are split into items with split) and applies
// the combinator stored for the longest matched prefix.
// It fails with errMessage when no key matches.
func MapTree[T comparable, P any, K comparable, V any](
	errMessage string,
	cases map[T]Combinator[K, P, V],
	split func(T) []K,
) Combinator[K, P, V] {
	tree := NewLongestPrefixTree(cases, split)

	var null V

	return func(buf Buffer[K, P]) (V, Error[P]) {
		pos := buf.Position()

		parse, err := tree.Lookup(buf)
		if err != nil {
			return null, NewParseError(pos, err.Error())
		}

		if parse != nil {
			return parse(buf)
		}

		return null, NewParseError(pos, errMessage)
	}
}
