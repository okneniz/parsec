package tml

import (
	"fmt"
	"strings"
)

// Type is a type expression: a nullary constructor (int, bool), a
// variable, a function arrow or a tuple. Variables carry their
// substitution link, so unification mutates the type graph in place.
type Type interface {
	String() string
}

// TCon is a nullary type constructor: int or bool.
type TCon struct {
	Name string
}

// TVar is a type variable. Link is nil while the variable is unbound
// and points to the unified type otherwise.
type TVar struct {
	Link Type
	id   int
}

// TFun is a function type.
type TFun struct {
	Arg, Res Type
}

// TTuple is a product type.
type TTuple struct {
	Items []Type
}

// The two builtin types: tiny ml has no type declarations, so these
// are the only constructors.
var (
	intType  = &TCon{Name: "int"}
	boolType = &TCon{Name: "bool"}
)

func (t *TCon) String() string {
	return t.Name
}

func (v *TVar) String() string {
	root := prune(v)
	if _, stillVar := root.(*TVar); !stillVar {
		return root.String()
	}

	return "'t" + fmt.Sprint(v.id)
}

func (t *TFun) String() string {
	return typeString(t, 0)
}

func (t *TTuple) String() string {
	return typeString(t, 0)
}

// typeString renders t with the usual ML conventions: arrows are
// right-associative, tuples bind tighter than arrows. Variables are
// named 'a, 'b, ... in order of appearance.
func typeString(t Type, prec int) string {
	names := map[*TVar]string{}

	assignNames(t, names, 0)

	return render(t, prec, names)
}

// assignNames numbers the variables of t left to right.
func assignNames(t Type, names map[*TVar]string, next int) int {
	switch ty := prune(t).(type) {
	case *TVar:
		if _, seen := names[ty]; !seen {
			names[ty] = "'" + string(rune('a'+next))
			next++
		}

	case *TFun:
		next = assignNames(ty.Arg, names, next)
		next = assignNames(ty.Res, names, next)

	case *TTuple:
		for _, item := range ty.Items {
			next = assignNames(item, names, next)
		}
	}

	return next
}

// render prints t assuming every variable has a name; prec is the
// binding strength of the context, so that argument tuples and left
// operands of arrows get parenthesized.
func render(t Type, prec int, names map[*TVar]string) string {
	switch ty := prune(t).(type) {
	case *TCon:
		return ty.Name
	case *TVar:
		return names[ty]
	case *TTuple:
		parts := make([]string, 0, len(ty.Items))

		for _, item := range ty.Items {
			parts = append(parts, render(item, 2, names))
		}

		out := strings.Join(parts, " * ")

		if prec > 2 {
			return "(" + out + ")"
		}

		return out
	case *TFun:
		out := render(ty.Arg, 2, names) + " -> " + render(ty.Res, 1, names)

		if prec > 1 {
			return "(" + out + ")"
		}

		return out
	}

	panic("unreachable type")
}

// prune follows the substitution links and compresses them, so that
// every variable on the path points to its current root.
func prune(t Type) Type {
	v, ok := t.(*TVar)
	if !ok || v.Link == nil {
		return t
	}

	v.Link = prune(v.Link)

	return v.Link
}

// unify makes a and b the same type, or fails with a description of
// the conflict.
func unify(a, b Type) error {
	a, b = prune(a), prune(b)

	av, aIsVar := a.(*TVar)
	bv, bIsVar := b.(*TVar)

	switch {
	case aIsVar && bIsVar && av == bv:
		return nil
	case aIsVar:
		if occurs(av, b) {
			return fmt.Errorf("type error: infinite type %s = %s", a, b)
		}

		av.Link = b

		return nil
	case bIsVar:
		if occurs(bv, a) {
			return fmt.Errorf("type error: infinite type %s = %s", a, b)
		}

		bv.Link = a

		return nil
	}

	switch at := a.(type) {
	case *TCon:
		bt, ok := b.(*TCon)
		if !ok || at.Name != bt.Name {
			return fmt.Errorf("type error: cannot unify %s with %s", a, b)
		}

		return nil
	case *TFun:
		bt, ok := b.(*TFun)
		if !ok {
			return fmt.Errorf("type error: cannot unify %s with %s", a, b)
		}

		if err := unify(at.Arg, bt.Arg); err != nil {
			return err
		}

		return unify(at.Res, bt.Res)
	case *TTuple:
		bt, ok := b.(*TTuple)
		if !ok || len(at.Items) != len(bt.Items) {
			return fmt.Errorf("type error: cannot unify %s with %s", a, b)
		}

		for i := range at.Items {
			if err := unify(at.Items[i], bt.Items[i]); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("type error: cannot unify %s with %s", a, b)
}

// occurs reports whether the variable appears inside t, which would
// make the unification circular.
func occurs(v *TVar, t Type) bool {
	switch ty := prune(t).(type) {
	case *TVar:
		return ty == v
	case *TFun:
		return occurs(v, ty.Arg) || occurs(v, ty.Res)
	case *TTuple:
		for _, item := range ty.Items {
			if occurs(v, item) {
				return true
			}
		}
	}

	return false
}

// scheme is a type scheme: the type with some of its variables
// universally quantified.
type scheme struct {
	vars []*TVar
	t    Type
}

// freeVars lists the unbound variables of t.
func freeVars(t Type) []*TVar {
	var out []*TVar

	seen := map[*TVar]struct{}{}

	var walk func(t Type)

	walk = func(t Type) {
		switch ty := prune(t).(type) {
		case *TVar:
			if _, dup := seen[ty]; !dup {
				seen[ty] = struct{}{}
				out = append(out, ty)
			}

		case *TFun:
			walk(ty.Arg)
			walk(ty.Res)

		case *TTuple:
			for _, item := range ty.Items {
				walk(item)
			}
		}
	}

	walk(t)

	return out
}

// generalize quantifies every variable of t which is not free in the
// environment — the heart of let-polymorphism.
func generalize(env map[string]*scheme, t Type) *scheme {
	envFree := map[*TVar]struct{}{}

	for _, s := range env {
		for _, v := range freeVars(s.t) {
			envFree[v] = struct{}{}
		}
	}

	var vars []*TVar

	for _, v := range freeVars(t) {
		if _, locked := envFree[v]; !locked {
			vars = append(vars, v)
		}
	}

	return &scheme{vars: vars, t: t}
}

// instantiate replaces the quantified variables of s with fresh ones.
func (c *typeChecker) instantiate(s *scheme) Type {
	if len(s.vars) == 0 {
		return s.t
	}

	mapping := map[Type]Type{}

	for _, v := range s.vars {
		mapping[v] = c.freshVar()
	}

	return copyType(s.t, mapping)
}

// copyType rebuilds t with the mapped variables replaced; variables
// outside the mapping are shared, because they are constrained by the
// environment.
func copyType(t Type, mapping map[Type]Type) Type {
	if mapped, ok := mapping[prune(t)]; ok {
		return mapped
	}

	switch ty := prune(t).(type) {
	case *TFun:
		return &TFun{
			Arg: copyType(ty.Arg, mapping),
			Res: copyType(ty.Res, mapping),
		}

	case *TTuple:
		items := make([]Type, 0, len(ty.Items))

		for _, item := range ty.Items {
			items = append(items, copyType(item, mapping))
		}

		return &TTuple{Items: items}
	}

	return t
}
