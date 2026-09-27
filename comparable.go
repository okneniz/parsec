package parsec

// Eq succeeds when the next item is equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func Eq[T comparable, P any, S any](
	errMessage string,
	t T,
) Combinator[T, P, T, S] {
	return Satisfy[T, P, S](errMessage, true, func(x T) bool {
		return t == x
	})
}

// NotEq succeeds when the next item is not equal to t and returns it.
// Greedy: consumes the item even on failure (see Try).
func NotEq[T comparable, P any, S any](
	errMessage string,
	t T,
) Combinator[T, P, T, S] {
	return Satisfy[T, P, S](errMessage, true, func(x T) bool {
		return t != x
	})
}

// OneOf succeeds when the next item is one of data and returns it.
// Greedy: consumes the item even on failure (see Try).
func OneOf[T comparable, P any, S any](
	errMessage string,
	data ...T,
) Combinator[T, P, T, S] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P, S](errMessage, true, func(x T) bool {
		_, exists := m[x]
		return exists
	})
}

// NoneOf succeeds when the next item is none of data and returns it.
// Greedy: consumes the item even on failure (see Try).
func NoneOf[T comparable, P any, S any](
	errMessage string,
	data ...T,
) Combinator[T, P, T, S] {
	m := make(map[T]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return Satisfy[T, P, S](errMessage, true, func(x T) bool {
		_, exists := m[x]
		return !exists
	})
}

// SequenceOf expects the next items to be equal to data in the same order
// and returns them as a slice.
func SequenceOf[T comparable, P any, S any](
	errMessage string,
	data ...T,
) Combinator[T, P, []T, S] {
	return func(_ S, buffer Buffer[T, P]) ([]T, Error[P]) {
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
func Map[T any, P any, K comparable, O any, S any](
	errMessage string,
	cases map[K]O,
	c Combinator[T, P, K, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()

		token, err := c(state, buffer)
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
func MapAs[T any, P any, K comparable, O any, S any](
	errMessage string,
	cases map[K]Combinator[T, P, O, S],
	comb Combinator[T, P, K, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()

		key, err := comb(state, buffer)
		if err != nil {
			return null, err
		}

		parseValue, exists := cases[key]
		if !exists {
			return null, NewParseError(pos, errMessage)
		}

		return parseValue(state, buffer)
	}
}

// MapTree matches the input against the keys of cases using a
// longest-prefix trie (keys are split into items with split) and applies
// the combinator stored for the longest matched prefix.
// It fails with errMessage when no key matches.
func MapTree[T comparable, P any, K comparable, O any, S any](
	errMessage string,
	cases map[K]Combinator[T, P, O, S],
	split func(K) []T,
) Combinator[T, P, O, S] {
	tree := NewLongestPrefixTree(cases, split)

	var null O

	return func(state S, buf Buffer[T, P]) (O, Error[P]) {
		pos := buf.Position()

		parse, err := tree.Lookup(state, buf)
		if err != nil {
			return null, NewParseError(pos, err.Error())
		}

		if parse != nil {
			return parse(state, buf)
		}

		return null, NewParseError(pos, errMessage)
	}
}
