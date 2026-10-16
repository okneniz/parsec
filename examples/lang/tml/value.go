package tml

import (
	"strconv"
	"strings"
)

// Value is a tiny ml runtime value. The literals and the tuples are
// the terms themselves; a function is a knot — see VRec — so every
// value carries the Expr interface and renders itself.
type Value = Expr

// VRec is a function value: the argument patterns with the body, and
// the row of knot values its free indices resolve against, with the
// depth they count from. The row holds the function itself in its
// shadow slot — the substitution core cannot tie that knot, the row
// does — which is what makes a declared fun, and the mutual recursion
// of a whole top level, self-referential.
type VRec struct {
	Args  []Pat
	Body  Expr
	Row   []Value
	Depth int
}

func (v *VRec) String() string {
	return "<fn>"
}

// negBuiltin is the builtin ~ the initial row holds at its bottom: a
// marker the application case of the interpreter recognizes.
type negBuiltin struct{}

func (negBuiltin) String() string {
	return "<fn>"
}

// renderValue renders a value the way the programs and the tests
// read it: canonical integers and booleans, comma-separated tuples,
// <fn> for functions.
func renderValue(v Value) string {
	switch value := v.(type) {
	case IntLit:
		return strconv.FormatInt(value.Value, 10)
	case Tuple:
		parts := make([]string, 0, len(value.Items))

		for _, item := range value.Items {
			parts = append(parts, renderValue(item))
		}

		return "(" + strings.Join(parts, ", ") + ")"
	}

	return v.String()
}

// appendRow copies row with one more value on top: the rows never
// share a backing array, so the row a knot captured stays its own.
func appendRow(row []Value, v Value) []Value {
	out := make([]Value, len(row), len(row)+1)
	copy(out, row)

	return append(out, v)
}
