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
func Satisfy[S any](
	errMessage string,
	greedy bool,
	f parsec.Condition[byte],
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Satisfy[byte, int, S](errMessage, greedy, f)
}

// Any reads and returns the next byte from the buffer.
// It fails only at the end of the buffer.
func Any[S any]() parsec.Combinator[byte, int, byte, S] {
	return parsec.Any[byte, int, S]()
}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
// The state is not rewound: a combinator that changes it must restore
// the change itself before failing.
func Try[O any, S any](c parsec.Combinator[byte, int, O, S]) parsec.Combinator[byte, int, O, S] {
	return parsec.Try[byte, int, O, S](c)
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[A any, B any, O any, S any](
	pre parsec.Combinator[byte, int, A, S],
	c parsec.Combinator[byte, int, O, S],
	suf parsec.Combinator[byte, int, B, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Between(pre, c, suf)
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF[S any]() parsec.Combinator[byte, int, bool, S] {
	return parsec.EOF[byte, int, S]()
}

// Cast parses the input with c and transforms the result with f.
// It fails when f returns an error.
func Cast[A any, O any, S any](
	c parsec.Combinator[byte, int, A, S],
	f func(A) (O, error),
) parsec.Combinator[byte, int, O, S] {
	return parsec.Cast(c, f)
}

// Const consumes nothing and returns the given value.
func Const[O any, S any](value O) parsec.Combinator[byte, int, O, S] {
	return parsec.Const[byte, int, O, S](value)
}

// Fail consumes nothing and always fails with the given message.
func Fail[O any, S any](errMessage string) parsec.Combinator[byte, int, O, S] {
	return parsec.Fail[byte, int, O, S](errMessage)
}

// Number is a constraint for any fixed-size binary number:
// an integer or a float.
type Number interface {
	parsec.Integer | parsec.Float
}

// ReadAs reads n bytes and decodes them as a single value of type O
// in the given byte order, for example a big-endian uint32.
// It fails with errMessage when the buffer ends before n bytes are read.
func ReadAs[O Number, S any](
	n int,
	errMessage string,
	order binary.ByteOrder,
) parsec.Combinator[byte, int, O, S] {
	anything := Any[S]()

	return func(state S, buffer parsec.Buffer[byte, int]) (O, parsec.Error[int]) {
		pos := buffer.Position()

		var result O

		input, err := Count(n, errMessage, anything)(state, buffer)
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
