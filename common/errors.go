package common

import (
	"errors"
	"fmt"
)

var (
	ErrEndOfFile   = errors.New("end of file")
	ErrOutOfBounds = errors.New("out of bounds")
)

type Error[T any] interface {
	error
	Position() T
	Previous() []Error[T]
}

type ParseError[T any] struct {
	message  string
	position T
	previous []Error[T]
}

var _ Error[int] = ParseError[int]{}

func (err ParseError[T]) Error() string {
	return fmt.Sprintf("Parse error at %v: %s", err.position, err.message)
}

func (err ParseError[T]) Position() T {
	return err.position
}

func (err ParseError[T]) Previous() []Error[T] {
	return err.previous
}

func NewParseError[T any](pos T, message string, previous ...Error[T]) ParseError[T] {
	return ParseError[T]{
		position: pos,
		message:  message,
		previous: previous,
	}
}
