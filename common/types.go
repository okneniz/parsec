package common

// Combinator is the type of a single parsing step: it reads items of
// type T from a [Buffer] with positions of type P and produces a value
// of type S or a parse [Error]. Combinators are values, so they can be
// stored, passed around and combined freely.
type Combinator[T any, P any, S any] func(Buffer[T, P]) (S, Error[P])

// Condition is a predicate over a single input item, see [Satisfy].
type Condition[T any] func(T) bool

// Composer combines the results of two consecutive combinators into
// one value, see [And].
type Composer[T any, S any, B any] func(T, S) B

// Anything is a [Condition] that accepts any item.
// Useful with the Satisfy combinator.
func Anything[T any](x T) bool { return true }

// Nothing is a [Condition] that rejects any item.
// Useful with the Satisfy combinator.
func Nothing[T any](x T) bool { return false }
