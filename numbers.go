package parsec

import "strconv"

// ParseInt parses digits with the parseDigits combinator and converts them
// to an integer of type O using strconv.ParseInt semantics with the given
// base and bitSize. It fails with errMessage when the digits are invalid
// or do not fit into O.
func ParseInt[T any, P any, O Integer, S any](
	errMessage string,
	base, bitSize int,
	parseDigits Combinator[T, P, string, S],
) Combinator[T, P, O, S] {
	var null O

	parse := Cast(parseDigits, func(digits string) (O, error) {
		num, err := strconv.ParseInt(digits, base, bitSize)
		if err != nil {
			return null, err
		}

		return O(num), nil
	})

	return func(state S, buf Buffer[T, P]) (O, Error[P]) {
		pos := buf.Position()

		num, err := parse(state, buf)
		if err != nil {
			return null, NewParseError(pos, errMessage, err)
		}

		return num, nil
	}
}
