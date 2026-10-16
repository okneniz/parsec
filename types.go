package parsec

// Combinator is the type of a single parsing step: it reads items of
// type T from a [Buffer] with positions of type P and produces a value
// of type O or a parse [Error]. The state of type S is shared mutable
// context: the entry point hands it to the first combinator and every
// combinator hands it to the next one. A combinator that changes the
// state must restore the change before returning an error — [Try]
// rewinds the buffer, not the state — so a failed alternative leaves
// the state as it found it. Combinators are values, so they can be
// stored, passed around and combined freely.
type Combinator[T any, P any, O any, S any] func(state S, buf Buffer[T, P]) (O, Error[P])

// Condition is a predicate over a single input item, see [Satisfy].
type Condition[T any] func(T) bool

// Composer combines the results of two consecutive combinators into
// one value, see [And].
type Composer[A any, B any, O any] func(A, B) O

// Anything is a [Condition] that accepts any item.
// Useful with the Satisfy combinator.

func Anything[T any](x T) bool {
	return true
}

// Nothing is a [Condition] that rejects any item.
// Useful with the Satisfy combinator.

func Nothing[T any](x T) bool {
	return false
}
