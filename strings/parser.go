package strings

import (
	"github.com/okneniz/parsec/common"
)

// Parse - parse text by c combinator.
func Parse[T any](
	data []rune,
	parse common.Combinator[rune, Position, T],
) (T, common.Error[Position]) {
	buf := Buffer(data)
	return common.Parse[rune, Position, T](buf, parse)
}

// ParseString - parse text by c combinator.
func ParseString[T any](
	str string,
	parse common.Combinator[rune, Position, T],
) (T, common.Error[Position]) {
	return Parse([]rune(str), parse)
}
