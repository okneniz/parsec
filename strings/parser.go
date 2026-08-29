package strings

import (
	"github.com/okneniz/parsec/common"
)

// Parse applies the c combinator to data.
func Parse[T any](
	data []rune,
	parse common.Combinator[rune, Position, T],
) (T, common.Error[Position]) {
	buf := Buffer(data)
	return common.Parse[rune, Position, T](buf, parse)
}

// ParseString converts str to runes and applies the c combinator to them.
func ParseString[T any](
	str string,
	parse common.Combinator[rune, Position, T],
) (T, common.Error[Position]) {
	return Parse([]rune(str), parse)
}
