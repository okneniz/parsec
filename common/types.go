package common

// Combinator - type of parse combinator.
type Combinator[T any, P any, S any] func(Buffer[T, P]) (S, error)

// Condition - condition function or predicate.
type Condition[T any] func(T) bool

// Compose - helper function to compose two values from combinators.
type Composer[T any, S any, B any] func(T, S) B

// Anything - return true anyway.
// Useful with Satisfy combinator.
func Anything[T any](x T) bool { return true }

// Nothing - return false anyway.
// Useful with Satisfy combinator.
func Nothing[T any](x T) bool { return false }
