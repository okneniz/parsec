package strings

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/parsec"
)

type counterState struct {
	seen []rune
}

func TestStateFlowsThroughCombinators(t *testing.T) {
	t.Parallel()

	read := func(state *counterState, buf parsec.Buffer[rune, Position]) (rune, parsec.Error[Position]) {
		r, err := buf.Read(true)
		if err != nil {
			return 0, parsec.NewParseError(buf.Position(), err.Error())
		}

		state.seen = append(state.seen, r)

		return r, nil
	}

	// And hands the same state to both parts; the order of the
	// collected runes records the order of the applications
	pair := And(
		read,
		read,
		func(a, b rune) [2]rune { return [2]rune{a, b} },
	)

	state := &counterState{}

	result, err := pair(state, Buffer([]rune("ab")))
	require.NoError(t, err)
	require.Equal(t, [2]rune{'a', 'b'}, result)
	require.Equal(t, []rune{'a', 'b'}, state.seen)
}

func TestTryRewindsBufferNotState(t *testing.T) {
	t.Parallel()

	var mutations []string

	toucher := func(state *counterState, buf parsec.Buffer[rune, Position]) (rune, parsec.Error[Position]) {
		mutations = append(mutations, "touch")

		return 0, parsec.NewParseError(buf.Position(), "always fails")
	}

	seeker := Try(toucher)

	state := &counterState{}

	_, err := seeker(state, Buffer([]rune("a")))
	require.Error(t, err)
	require.Len(t, mutations, 1, "the combinator ran once")
	require.NotNil(t, state, "the same state came back")
}

func TestStatefulChoiceRestoresItsOwnState(t *testing.T) {
	t.Parallel()

	// the pushing combinator fails after changing the state and must
	// restore the change itself: the next alternative, and a retry of
	// the whole Choice, see the state as it was
	pushAndFail := func(state *counterState, buf parsec.Buffer[rune, Position]) (rune, parsec.Error[Position]) {
		state.seen = append(state.seen, 'x')
		defer func() {
			state.seen = state.seen[:len(state.seen)-1]
		}()

		return 0, parsec.NewParseError(buf.Position(), "fails after push")
	}

	fallback := Eq[*counterState]("expected 'a'", 'a')

	alternatives := Choice("expected one of the alternatives",
		Try(pushAndFail),
		Try(fallback),
	)

	state := &counterState{}

	result, err := alternatives(state, Buffer([]rune("a")))
	require.NoError(t, err)
	require.Equal(t, 'a', result)
	require.Empty(t, state.seen)
}
