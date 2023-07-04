package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Parse[T any](
	data []rune,
	parse p.Combinator[rune, Position, T],
) (T, error) {
	buf := Buffer(data)
	return p.Parse[rune, Position, T](buf, parse)
}

func ParseString[T any](
	str string,
	parse p.Combinator[rune, Position, T],
) (T, error) {
	return Parse([]rune(str), parse)
}
