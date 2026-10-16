package strings

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/parsec"
)

type testState struct {
	seen    []rune
	counter int
}

func (c *testState) Clone() *testState {
	cln := testState{
		seen:    make([]rune, 0, len(c.seen)),
		counter: c.counter,
	}

	cln.seen = append(cln.seen, c.seen...)

	return &cln
}

func (c *testState) Recover(cln *testState) {
	c.seen = cln.seen
	c.counter = cln.counter
}

func TestStateFlowsThroughCombinators(t *testing.T) {
	t.Parallel()

	read := func(
		state *testState,
		buf parsec.Buffer[rune, Position],
	) (rune, parsec.Error[Position]) {
		r, err := buf.Read(true)
		if err != nil {
			return 0, parsec.NewParseError(buf.Position(), err.Error())
		}

		state.counter++
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

	state := &testState{}

	result, err := pair(state, Buffer([]rune("ab")))
	require.NoError(t, err)
	require.Equal(t, [2]rune{'a', 'b'}, result)
	require.Equal(t, []rune{'a', 'b'}, state.seen)
	require.Equal(t, 2, state.counter)
}

func TestTryRewindsStateToo(t *testing.T) {
	t.Parallel()

	var mutations []string

	toucher := func(
		state *testState,
		buf parsec.Buffer[rune, Position],
	) (rune, parsec.Error[Position]) {
		state.counter = 100
		mutations = append(mutations, "touch")
		return 0, parsec.NewParseError(buf.Position(), "always fails")
	}

	seeker := Try(toucher)
	state := &testState{}

	_, err := seeker(state, Buffer([]rune("a")))
	require.Error(t, err)
	require.Equal(t, 0, state.counter)
	require.Len(t, mutations, 1, "the combinator ran once")
	require.NotNil(t, state, "the same state came back")
}

func TestStatefulChoiceRestoresItsOwnState(t *testing.T) {
	t.Parallel()

	// the pushing combinator fails after changing the state and must
	// restore the change itself: the next alternative, and a retry of
	// the whole Choice, see the state as it was
	pushAndFail := func(
		state *testState,
		buf parsec.Buffer[rune, Position],
	) (rune, parsec.Error[Position]) {
		state.seen = append(state.seen, 'x')
		defer func() {
			state.seen = state.seen[:len(state.seen)-1]
		}()

		return 0, parsec.NewParseError(buf.Position(), "fails after push")
	}

	fallback := Eq[*testState]("expected 'a'", 'a')

	alternatives := Choice(
		"expected one of the alternatives",
		Try(pushAndFail),
		Try(fallback),
	)

	state := &testState{}

	result, err := alternatives(state, Buffer([]rune("a")))
	require.NoError(t, err)
	require.Equal(t, 'a', result)
	require.Empty(t, state.seen)
}
