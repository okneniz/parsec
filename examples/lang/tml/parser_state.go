package tml

// state is the parser state of tiny ml: the stack of the binder
// names a nested expression may reference, innermost last. The
// builtin ~ sits at the bottom of the initial stack; every fn
// pattern, every declaration and every let pushes its binders for
// the nested scopes and pops them on the way out — the state half
// of the backtracking contract, because Try rewinds only the buffer.
type state struct {
	binders []string
}

func newState() *state {
	return &state{binders: []string{"~"}}
}

func (s *state) Clone() *state {
	clone := state{
		binders: make([]string, 0, len(s.binders)),
	}

	clone.binders = append(clone.binders, s.binders...)

	return &clone
}

func (s *state) Recover(clone *state) {
	s.binders = clone.binders
}

func (s *state) push(names ...string) {
	s.binders = append(s.binders, names...)
}

// popTo truncates the stack back to n binders.
func (s *state) popTo(n int) {
	s.binders = s.binders[:n]
}

// lookup returns the de Bruijn index of the innermost binder named
// name: 0 is the innermost binder of all.
func (s *state) lookup(name string) (int, bool) {
	for i := len(s.binders) - 1; i >= 0; i-- {
		if s.binders[i] == name {
			return len(s.binders) - 1 - i, true
		}
	}

	return 0, false
}

// builtinVar is the reference to a builtin — the binder at the
// bottom of the stack, where ~ sits alone.
func (s *state) builtinVar() Var {
	return Var{Index: len(s.binders) - 1}
}

// depth returns the number of binders on the stack.
func (s *state) depth() int {
	return len(s.binders)
}
