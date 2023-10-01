package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Range - succeeds for any item which include in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Range(from rune, to rune) p.Combinator[rune, Position, rune] {
	return p.Range[rune, Position](from, to)
}

// NotRange - succeeds for any item which not included in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotRange(from rune, to rune) p.Combinator[rune, Position, rune] {
	return p.NotRange[rune, Position](from, to)
}

// Gt - succeeds for any item which greater than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gt(t rune) p.Combinator[rune, Position, rune] {
	return Satisfy(true, func(x rune) bool {
		return x > t
	})
}

// Gte - succeeds for any item which greater than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gte(t rune) p.Combinator[rune, Position, rune] {
	return Satisfy(true, func(x rune) bool {
		return x >= t
	})
}

// Lt - succeeds for any item which less than input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lt(t rune) p.Combinator[rune, Position, rune] {
	return Satisfy(true, func(x rune) bool {
		return x < t
	})
}

// Lte - succeeds for any item which less than or equal input value.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lte(t rune) p.Combinator[rune, Position, rune] {
	return Satisfy(true, func(x rune) bool {
		return x <= t
	})
}
