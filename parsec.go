package parsec

func Parse[T any, S any](buffer Buffer[T], parse Combinator[T, S]) (S, error) {
	result, err := parse(buffer)
	if err != nil {
		return *new(S), err
	}

	return result, nil
}

func ParseBytes[S any](data []byte, parse Combinator[byte, S]) (S, error) {
	buf := BytesBuffer(data)
	return Parse[byte, S](buf, parse)
}
