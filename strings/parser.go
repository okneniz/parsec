package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Parse - parse text by c combinator.
func Parse[T any](
	data []rune,
	parse p.Combinator[rune, Position, T],
) (T, error) {
	buf := Buffer(data)
	return p.Parse[rune, Position, T](buf, parse)
}

// ParseString - parse text by c combinator.
func ParseString[T any](
	str string,
	parse p.Combinator[rune, Position, T],
) (T, error) {
	return Parse([]rune(str), parse)
}
