package common

// Eq - succeeds for any item which equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Eq[T comparable, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return t == x
	})
}

// NotEq - succeeds for any item which not equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotEq[T comparable, P any](
	errMessage string,
	t T,
) Combinator[T, P, T] {
	return Satisfy[T, P](errMessage, true, func(x T) bool {
		return t != x
	})
}

// OneOf - succeeds for any item which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
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

// NoneOf - succeeds for any item which not included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
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

// SequenceOf - expects a sequence of elements in the buffer
// equal to the input data sequence. If expectations are not met,
// returns ParseError error.
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

// Map - Reads one element from the input buffer using the combinator,
// then uses the resulting element to obtain a value from the map cases and try to
// match it in cases map passed by first argument.
// If the value is not found then it returns ParseError error.
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

// MapAs - Read one element from the input buffer using the combinator,
// then match the resulting item to obtain a value from map cases and try to match it
// in cases map passed by first argument.
// If the value it not found then it returns ParseError error.
// Otherwise try to parse input data by combinator from cases.
func MapAs[T any, P any, K comparable, V any](
	errMessage string,
	cases map[K]Combinator[T, P, V],
	comb Combinator[T, P, K],
) Combinator[T, P, V] {
	var null V

	// TODO : make error message

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

// MapTree - Reads element from the input buffer using the combinator and
// match it in on the fly by cases map passed by second argument.
// Try to parse longest prefix.
// If the value is not found then it returns ParseError error.
// This combinator use special trie-like structure for text matching.
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
