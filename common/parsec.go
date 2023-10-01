package common

// Parse - parse data from buffer by c combinator.
func Parse[T any, P any, S any](
	buffer Buffer[T, P],
	c Combinator[T, P, S],
) (S, error) {
	result, err := c(buffer)
	if err != nil {
		return *new(S), err
	}

	return result, nil
}
