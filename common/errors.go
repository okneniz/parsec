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
}

type ParseError[T any] struct {
	message  string
	position T
}

var _ Error[int] = ParseError[int]{}

func (err ParseError[T]) Error() string {
	return fmt.Sprintf("Parse error at %v: %s", err.position, err.message)
}

func (err ParseError[T]) Position() T {
	return err.position
}

func NewParseError[T any](pos T, message string) ParseError[T] {
	return ParseError[T]{
		position: pos,
		message:  message,
	}
}
