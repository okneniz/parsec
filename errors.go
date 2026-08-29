package parsec

import (
	"errors"
	"fmt"
)

// ErrEndOfFile is returned by a Buffer when reading past the end of input.
var ErrEndOfFile = errors.New("end of file")

// ErrOutOfBounds is returned by a Buffer when seeking to an invalid position.
var ErrOutOfBounds = errors.New("out of bounds")

// Error is a parse error bound to an input position.
// Nested combinators accumulate their failures in Previous,
// which makes the error chain of a failed Choice readable.
type Error[T any] interface {
	error
	Position() T
	Previous() []Error[T]
}

// ParseError is the default Error implementation.
type ParseError[T any] struct {
	message  string
	position T
	previous []Error[T]
}

var _ Error[int] = ParseError[int]{}

// Error implements the error interface.
func (err ParseError[T]) Error() string {
	return fmt.Sprintf("Parse error at %v: %s", err.position, err.message)
}

// Position returns the input position the error is bound to.
func (err ParseError[T]) Position() T {
	return err.position
}

// Previous returns the errors of the nested combinators that led to this one.
func (err ParseError[T]) Previous() []Error[T] {
	return err.previous
}

// NewParseError creates a ParseError bound to pos, with optional
// previous errors attached.
func NewParseError[T any](pos T, message string, previous ...Error[T]) ParseError[T] {
	return ParseError[T]{
		position: pos,
		message:  message,
		previous: previous,
	}
}
