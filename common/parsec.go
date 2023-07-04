package common

func Parse[T any, P any, S any](
	buffer Buffer[T, P],
	parse Combinator[T, P, S],
) (S, error) {
	result, err := parse(buffer)
	if err != nil {
		return *new(S), err
	}

	return result, nil
}
