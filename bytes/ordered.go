package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Range - succeeds for any bytes which include in input range.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Range(
	errMessage string,
	from, to byte,
) common.Combinator[byte, int, byte] {
	return common.Range[byte, int](errMessage, from, to)
}

// NotRange - succeeds for any byte which not included in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotRange(
	errMessage string,
	from, to byte,
) common.Combinator[byte, int, byte] {
	return common.NotRange[byte, int](errMessage, from, to)
}

// Gt - succeeds for any byte which greater than input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gt(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return Satisfy(errMessage, true, func(x byte) bool {
		return x > t
	})
}

// Gte - succeeds for any byte which greater than or equal input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gte(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return Satisfy(errMessage, true, func(x byte) bool {
		return x >= t
	})
}

// Lt - succeeds for any byte which less than input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lt(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return Satisfy(errMessage, true, func(x byte) bool {
		return x < t
	})
}

// Lte - succeeds for any byte which less than or equal input byte.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lte(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return Satisfy(errMessage, true, func(x byte) bool {
		return x <= t
	})
}
