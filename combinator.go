package parsec

import (
	"errors"
	"fmt"
)

// Satisfy succeeds for any item for which the supplied function f returns true
// and returns that item.
// If greedy is true, buffer keeps position after reading - even when f fails,
// so the failed combinator consumes exactly one item.
// Wrap it in Try if you need the position to be restored on failure.
func Satisfy[T any, P any, S any](
	errMessage string,
	greedy bool,
	f Condition[T],
) Combinator[T, P, T, S] {
	var null T

	return func(_ S, buffer Buffer[T, P]) (T, Error[P]) {
		pos := buffer.Position()

		token, err := buffer.Read(greedy)
		if err != nil {
			return null, NewParseError(pos, errMessage)
		}

		if f(token) {
			return token, nil
		}

		return null, NewParseError(pos, errMessage)
	}
}

// Any reads and returns the next item from the buffer.
// It fails only at the end of the buffer.
func Any[T any, P any, S any]() Combinator[T, P, T, S] {
	var null T

	return func(_ S, buffer Buffer[T, P]) (T, Error[P]) {
		pos := buffer.Position()

		token, err := buffer.Read(true)
		if err != nil {
			return null, NewParseError(pos, err.Error())
		}

		return token, nil
	}
}

// Between parses open, then body, then close,
// and returns only the result of body.
func Between[T any, P any, A any, B any, O any, S any](
	pre Combinator[T, P, A, S],
	c Combinator[T, P, O, S],
	suf Combinator[T, P, B, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		_, err := pre(state, buffer)
		if err != nil {
			return null, err
		}

		body, err := c(state, buffer)
		if err != nil {
			return null, err
		}

		_, err = suf(state, buffer)
		if err != nil {
			return null, err
		}

		return body, nil
	}
}

// EOF reports whether the buffer is fully consumed.
// It never fails and consumes nothing.
func EOF[T any, P any, S any]() Combinator[T, P, bool, S] {
	return func(_ S, buffer Buffer[T, P]) (bool, Error[P]) {
		if buffer.IsEOF() {
			return true, nil
		}

		return false, nil
	}
}

// Cast parses the input with c and transforms the result with f.
// It fails with errMessage when f returns an error.
func Cast[T any, P any, A any, O any, S any](
	c Combinator[T, P, A, S],
	cast func(A) (O, error),
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()

		result, err := c(state, buffer)
		if err != nil {
			return null, err
		}

		value, castError := cast(result)
		if castError != nil {
			return value, NewParseError(pos, castError.Error())
		}

		return value, nil
	}
}

// Unexpected builds a combinator which fails when the step succeeds,
// describing the item the step produced: a negative lookahead. It is
// the natural shape of "nothing else may parse here" checks — the
// end-of-input check of a parser feeds its own stepper as the step
// and gets "unexpected <item>, expected end of input" for leftovers.
//
// A positional error of the step is returned untouched; the bare
// ErrEndOfFile becomes "unexpected end of input".
func Unexpected[T any, P any, O any, S any](
	step func(state S, buf Buffer[T, P]) (O, error),
	context string,
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buf Buffer[T, P]) (O, Error[P]) {
		pos := buf.Position()

		result, err := step(state, buf)
		if err != nil {
			if parseErr, ok := errors.AsType[Error[P]](err); ok {
				return null, parseErr
			}

			if errors.Is(err, ErrEndOfFile) {
				return null, NewParseError(pos, "unexpected end of input, "+context)
			}

			return null, NewParseError(pos, err.Error())
		}

		return null, NewParseError(pos, fmt.Sprintf("unexpected %v, %s", result, context))
	}
}

// Const consumes nothing and returns the given value.
func Const[T any, P any, O any, S any](value O) Combinator[T, P, O, S] {
	return func(_ S, _ Buffer[T, P]) (O, Error[P]) {
		return value, nil
	}
}

// Fail consumes nothing and always fails with the given message.
func Fail[T any, P any, O any, S any](errMessage string) Combinator[T, P, O, S] {
	var null O

	return func(_ S, buffer Buffer[T, P]) (O, Error[P]) {
		return null, NewParseError(buffer.Position(), errMessage)
	}
}
