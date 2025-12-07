package common

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func ParseInt[T any, P any, S constraints.Integer](
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
