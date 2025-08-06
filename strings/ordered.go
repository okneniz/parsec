package strings

import (
	"github.com/okneniz/parsec/common"
)

// Range - succeeds for any item which include in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Range(
	errMessage string,
	from, to rune,
) common.Combinator[rune, Position, rune] {
	return common.Range[rune, Position](errMessage, from, to)
}

// NotRange - succeeds for any item which not included in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotRange(
	errMessage string,
	from, to rune,
) common.Combinator[rune, Position, rune] {
	return common.NotRange[rune, Position](errMessage, from, to)
}

// Gt - succeeds for any item which greater than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gt(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return Satisfy(errMessage, true, func(x rune) bool {
		return x > t
	})
}

// Gte - succeeds for any item which greater than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gte(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return Satisfy(errMessage, true, func(x rune) bool {
		return x >= t
	})
}

// Lt - succeeds for any item which less than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lt(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return Satisfy(errMessage, true, func(x rune) bool {
		return x < t
	})
}

// Lte - succeeds for any item which less than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lte(
	errMessage string,
	t rune,
) common.Combinator[rune, Position, rune] {
	return Satisfy(errMessage, true, func(x rune) bool {
		return x <= t
	})
}
