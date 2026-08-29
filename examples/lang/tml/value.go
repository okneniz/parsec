package tml

import (
	"fmt"
	"sort"
	"strings"
)

// Value is a tiny ml runtime value: an integer, a boolean, a tuple or
// a closure.
type Value interface {
	String() string
}

// Env is an environment of bindings; the closures capture it by
// reference, which is what makes recursion work.
type Env map[string]Value

// VInt is an integer value.
type VInt int64

// VBool is a boolean value.
type VBool bool

// VTuple is a tuple value.
type VTuple []Value

// VClosure is a function value: one argument pattern, the body, and
// the environment captured by reference.
type VClosure struct {
	Arg  Pat
	Body Expr
	Env  Env
}

func (v VInt) String() string {
	return fmt.Sprint(int64(v))
}

func (v VBool) String() string {
	return fmt.Sprint(bool(v))
}

func (v VTuple) String() string {
	out := "("

	for i, item := range v {
		if i > 0 {
			out += ", "
		}

		out += item.String()
	}

	return out + ")"
}

func (v *VClosure) String() string {
	return "<fn>"
}

// String renders the bindings one per line, in name order — the
// stable format the tests and the examples check against.
func (env Env) String() string {
	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}

	sort.Strings(names)

	lines := make([]string, 0, len(env))
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("%s = %s", name, env[name]))
	}

	return strings.Join(lines, "\n")
}
