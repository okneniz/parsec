package parsec

type Recoverable[T any] interface {
	Clone() T
	Recover(T)
}

// Stateless is the state of a parser that keeps none: instantiate
// the state parameter S with the type and start the parse with the
// value Stateless{}.
type Stateless struct{}

func (s Stateless) Clone() Stateless {
	return s
}

func (Stateless) Recover(_ Stateless) {}

// Try applies c and, when it fails, rewinds the buffer to the previous
// position and the state to the snapshot Clone took before c ran.
// This is the backtracking primitive of the library: greedy combinators
// like Satisfy, Eq, Range and others consume input even when they fail,
// so every combinator which can fail inside an alternative branch
// (see Choice, Or) or in a loop (see Many and similar) must be wrapped
// in Try to make backtracking possible. The state is rewound on
// failure only: the changes of a successful application stay.
func Try[T any, P any, O any, S Recoverable[S]](
	c Combinator[T, P, O, S],
) Combinator[T, P, O, S] {
	var null O

	return func(state S, buffer Buffer[T, P]) (O, Error[P]) {
		pos := buffer.Position()
		clone := state.Clone()

		result, err := c(state, buffer)
		if err != nil {
			state.Recover(clone)

			if seekErr := buffer.Seek(pos); seekErr != nil {
				return null, NewParseError(
					buffer.Position(),
					seekErr.Error(),
					err,
				)
			}

			return null, err
		}

		return result, nil
	}
}
