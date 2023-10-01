package bytes

import (
	"bytes"
	"encoding/binary"
	p "github.com/okneniz/parsec/common"
	"golang.org/x/exp/constraints"
)

// Satisfy - succeeds for any byte for which the supplied function f returns true.
// Returns the byte that is actually readed from input buffer.
// if greedy buffer keep position after reading.
func Satisfy(
	greedy bool,
	f p.Condition[byte],
) p.Combinator[byte, int, byte] {
	return p.Satisfy[byte, int](greedy, f)
}

// Any - returns the readed byte.
func Any() p.Combinator[byte, int, byte] {
	return p.Any[byte, int]()
}

// Try - try to use c combinator, if it falls, it returns buffer to the previous position.
func Try[T any](c p.Combinator[byte, int, T]) p.Combinator[byte, int, T] {
	return p.Try[byte, int, T](c)
}

// Between - parse sequence of input combinators, skip first and last results.
func Between[T any, S any, B any](
	pre p.Combinator[byte, int, T],
	c p.Combinator[byte, int, S],
	suf p.Combinator[byte, int, B],
) p.Combinator[byte, int, S] {
	return p.Between(pre, c, suf)
}

// EOF - checks that buffer reading has finished.
func EOF() p.Combinator[byte, int, bool] {
	return p.EOF[byte, int]()
}

// Cast - parse data by c combinator and apply to f function.
// Return result of f function.
func Cast[T any, S any](
	c p.Combinator[byte, int, T],
	f func(T) (S, error),
) p.Combinator[byte, int, S] {
	return p.Cast(c, f)
}

// Const - doesn't read anything, just return the input value.
func Const[S any](value S) p.Combinator[byte, int, S] {
	return p.Const[byte, int, S](value)
}

// Fail - doesn't read anything, just return input error.
func Fail[S any](err error) p.Combinator[byte, int, S] {
	return p.Fail[byte, int, S](err)
}

type Number interface {
	constraints.Integer | constraints.Float
}

// ReadAs - read n of bytes,
// decode it in binary order passed by second argument
// and return it.
func ReadAs[T Number](
	n int,
	order binary.ByteOrder,
) p.Combinator[byte, int, T] {
	return func(buffer p.Buffer[byte, int]) (T, error) {
		var result T

		input, err := Count(n, Any())(buffer)
		if err != nil {
			return result, err
		}

		buf := bytes.NewReader(input)

		err = binary.Read(buf, order, &result)
		if err != nil {
			return result, err
		}

		return result, nil
	}
}
