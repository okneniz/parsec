package bytes

import (
	"bytes"
	"encoding/binary"

	"github.com/okneniz/parsec/common"
	"golang.org/x/exp/constraints"
)

// Satisfy - succeeds for any byte for which the supplied function f returns true.
// Returns the byte that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy(
	errMessage string,
	greedy bool,
	f common.Condition[byte],
) common.Combinator[byte, int, byte] {
	return common.Satisfy[byte, int](errMessage, greedy, f)
}

// Any - returns the readed byte.
func Any() common.Combinator[byte, int, byte] {
	return common.Any[byte, int]()
}

// Try - try to use c combinator, if it falls, it returns buffer to the previous position.
func Try[T any](c common.Combinator[byte, int, T]) common.Combinator[byte, int, T] {
	return common.Try[byte, int, T](c)
}

// Between - parse sequence of input combinators, skip first and last results.
func Between[T any, S any, B any](
	pre common.Combinator[byte, int, T],
	c common.Combinator[byte, int, S],
	suf common.Combinator[byte, int, B],
) common.Combinator[byte, int, S] {
	return common.Between(pre, c, suf)
}

// EOF - checks that buffer reading has finished.
func EOF() common.Combinator[byte, int, bool] {
	return common.EOF[byte, int]()
}

// Cast - parse data by c combinator and apply to f function.
// Return result of f function.
func Cast[T any, S any](
	c common.Combinator[byte, int, T],
	f func(T) (S, error),
) common.Combinator[byte, int, S] {
	return common.Cast(c, f)
}

// Const - doesn't read anything, just return the input value.
func Const[S any](value S) common.Combinator[byte, int, S] {
	return common.Const[byte, int, S](value)
}

// Fail - doesn't read anything, just return input error.
func Fail[S any](errMessage string) common.Combinator[byte, int, S] {
	return common.Fail[byte, int, S](errMessage)
}

type Number interface {
	constraints.Integer | constraints.Float
}

// ReadAs - read n of bytes,
// decode it in binary order passed by second argument
// and return it.
func ReadAs[T Number](
	n int,
	errMessage string,
	order binary.ByteOrder,
) common.Combinator[byte, int, T] {
	anything := Any()

	return func(buffer common.Buffer[byte, int]) (T, common.Error[int]) {
		pos := buffer.Position()

		var result T

		input, err := Count(n, errMessage, anything)(buffer)
		if err != nil {
			return result, err
		}

		buf := bytes.NewReader(input)

		readErr := binary.Read(buf, order, &result)
		if readErr != nil {
			return result, common.NewParseError(pos, readErr.Error())
		}

		return result, nil
	}
}
