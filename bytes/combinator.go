package bytes

import (
	"bytes"
	"encoding/binary"

	"github.com/okneniz/parsec"
)

// Satisfy succeeds for any byte for which the supplied function f returns
// true and returns that byte. If greedy is true, the buffer keeps the
// position after reading - even when f fails, so the failed combinator
// consumes exactly one byte. Wrap it in Try if you need the position
// to be restored on failure.
func Satisfy(
	errMessage string,
	greedy bool,
	f parsec.Condition[byte],
) parsec.Combinator[byte, int, byte] {
	return parsec.Satisfy[byte, int](errMessage, greedy, f)
}

// Any reads and returns the next byte from the buffer.
// It fails only at the end of the buffer.
func Any() parsec.Combinator[byte, int, byte] {
	return parsec.Any[byte, int]()
}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
func Try[T any](c parsec.Combinator[byte, int, T]) parsec.Combinator[byte, int, T] {
	return parsec.Try[byte, int, T](c)
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[T any, S any, B any](
	pre parsec.Combinator[byte, int, T],
	c parsec.Combinator[byte, int, S],
	suf parsec.Combinator[byte, int, B],
) parsec.Combinator[byte, int, S] {
	return parsec.Between(pre, c, suf)
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF() parsec.Combinator[byte, int, bool] {
	return parsec.EOF[byte, int]()
}

// Cast parses the input with c and transforms the result with f.
// It fails when f returns an error.
func Cast[T any, S any](
	c parsec.Combinator[byte, int, T],
	f func(T) (S, error),
) parsec.Combinator[byte, int, S] {
	return parsec.Cast(c, f)
}

// Const consumes nothing and returns the given value.
func Const[S any](value S) parsec.Combinator[byte, int, S] {
	return parsec.Const[byte, int, S](value)
}

// Fail consumes nothing and always fails with the given message.
func Fail[S any](errMessage string) parsec.Combinator[byte, int, S] {
	return parsec.Fail[byte, int, S](errMessage)
}

// Number is a constraint for any fixed-size binary number:
// an integer or a float.
type Number interface {
	parsec.Integer | parsec.Float
}

// ReadAs reads n bytes and decodes them as a single value of type T
// in the given byte order, for example a big-endian uint32.
// It fails with errMessage when the buffer ends before n bytes are read.
func ReadAs[T Number](
	n int,
	errMessage string,
	order binary.ByteOrder,
) parsec.Combinator[byte, int, T] {
	anything := Any()

	return func(buffer parsec.Buffer[byte, int]) (T, parsec.Error[int]) {
		pos := buffer.Position()

		var result T

		input, err := Count(n, errMessage, anything)(buffer)
		if err != nil {
			return result, err
		}

		buf := bytes.NewReader(input)

		readErr := binary.Read(buf, order, &result)
		if readErr != nil {
			return result, parsec.NewParseError(pos, readErr.Error())
		}

		return result, nil
	}
}
