package common

type Combinator[T any, P any, S any] func(Buffer[T, P]) (S, error)

type Condition[T any] func(T) bool
type Composer[T any, S any, B any] func(T, S) B

func Anything[T any](x T) bool { return true }
func Nothing[T any](x T) bool  { return false }
