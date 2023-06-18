package parsec

func Parse[T any, S any](buffer Buffer[T], parse Combinator[T, S]) (S, bool) {
	result, ok := parse(buffer)
	if ok {
		return result, ok
	}

	return *new(S), false
}

func ParseBytes[S any](data []byte, parse Combinator[byte, S]) (S, bool) {
	buf := BytesBuffer(data)
	return Parse[byte, S](buf, parse)
}
