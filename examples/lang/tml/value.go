package tml

import (
	"fmt"
)

// Value is a tiny ml runtime value: an integer, a boolean, a tuple or
// a closure.
type Value interface {
	String() string
}

// cell is one binder's slot. The top-level frame is a row of cells
// the declarations fill as the program runs, so a closure built
// earlier sees a function declared later — the mutual recursion of
// tiny ml.
type cell struct {
	val Value
}

// VInt is an integer value.
type VInt int64

// VBool is a boolean value.
type VBool bool

// VTuple is a tuple value.
type VTuple []Value

// VClosure is a function value: one argument pattern, the body, the
// program frame it reads the top level through, and the environment
// snapshot it captured.
type VClosure struct {
	Arg   Pat
	Body  Expr
	Frame []*cell
	Env   []*cell
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
