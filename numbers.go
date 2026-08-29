package parsec

import "strconv"

// ParseInt parses digits with the parseDigits combinator and converts them
// to an integer of type S using strconv.ParseInt semantics with the given
// base and bitSize. It fails with errMessage when the digits are invalid
// or do not fit into S.
func ParseInt[T any, P any, S Integer](
	errMessage string,
	base, bitSize int,
	parseDigits Combinator[T, P, string],
) Combinator[T, P, S] {
	var null S

	parse := Cast(parseDigits, func(digits string) (S, error) {
		num, err := strconv.ParseInt(digits, base, bitSize)
		if err != nil {
			return null, err
		}

		return S(num), nil
	})

	return func(buf Buffer[T, P]) (S, Error[P]) {
		pos := buf.Position()

		num, err := parse(buf)
		if err != nil {
			return null, NewParseError(pos, errMessage, err)
		}

		return num, nil
	}
}
