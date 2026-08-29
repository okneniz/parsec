package common

// Parse applies the c combinator to the buffer and returns its result.
// It is a small convenience wrapper: c(buffer) does the same.
func Parse[T any, P any, S any](
	buffer Buffer[T, P],
	c Combinator[T, P, S],
) (S, Error[P]) {
	result, err := c(buffer)
	if err != nil {
		return result, err
	}

	return result, nil
}
