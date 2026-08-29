package tml

import (
	"fmt"
)

// typeChecker is one run of Hindley-Milner inference: a counter for fresh
// type variables. The substitution lives in the variables themselves.
type typeChecker struct {
	fresh int
}

// freshVar creates a new unbound type variable.
func (c *typeChecker) freshVar() *TVar {
	c.fresh++

	return &TVar{id: c.fresh}
}

// Infer parses and typechecks src with Hindley-Milner inference and
// returns the types of the top-level bindings. Top-level functions
// see each other, so they may be mutually recursive.
func Infer(src string) (map[string]Type, error) {
	decls, err := Parse(src)
	if err != nil {
		return nil, err
	}

	return inferProgram(decls)
}

// inferProgram typechecks already parsed declarations.
func inferProgram(decls []Decl) (map[string]Type, error) {
	c := &typeChecker{}
	env := map[string]*scheme{}

	// pre-bind every top-level function with a fresh variable, so
	// that a body may reference functions declared later; the shared
	// variable is what carries the mutual constraints
	for _, d := range decls {
		fun, ok := d.(FunDecl)
		if !ok {
			continue
		}

		env[fun.Name] = &scheme{t: c.freshVar()}
	}

	for _, d := range decls {
		if err := c.inferDecl(env, d, true); err != nil {
			return nil, err
		}
	}

	out := map[string]Type{}

	for name, s := range env {
		out[name] = prune(s.t)
	}

	return out, nil
}

// without copies env minus one binding.
func without(env map[string]*scheme, name string) map[string]*scheme {
	out := make(map[string]*scheme, len(env))

	for k, v := range env {
		if k != name {
			out[k] = v
		}
	}

	return out
}

// inferDecl types one declaration, extending env in place. At the top
// level a function name is already pre-bound; inside a let it is
// fresh, shadowing any outer binding.
func (c *typeChecker) inferDecl(env map[string]*scheme, d Decl, top bool) error {
	switch decl := d.(type) {
	case ValDecl:
		for _, b := range decl.Binds {
			bindings := map[string]Type{}

			patType := c.patternType(b.Pat, bindings)

			valueType, err := c.inferExpr(env, b.E)
			if err != nil {
				return err
			}

			if err := unify(patType, valueType); err != nil {
				return err
			}

			// let-polymorphism: each variable of the pattern is
			// generalized over what the environment does not lock
			for name, t := range bindings {
				env[name] = generalize(env, t)
			}
		}

		return nil
	case FunDecl:
		self := &scheme{t: c.freshVar()}

		if pre, bound := env[decl.Name]; bound && top {
			self = pre
		}

		// the parameters live in an inner scope: they must not leak
		// into the environment the declaration extends
		inner := map[string]*scheme{decl.Name: self}

		args := make([]Type, 0, len(decl.Args))

		for _, arg := range decl.Args {
			bindings := map[string]Type{}

			t := c.patternType(arg, bindings)

			for name, bt := range bindings {
				inner[name] = &scheme{t: bt}
			}

			args = append(args, t)
		}

		bodyType, err := c.inferExpr(extend(env, inner), decl.Body)
		if err != nil {
			return err
		}

		funType := bodyType

		for i := len(args) - 1; i >= 0; i-- {
			funType = &TFun{Arg: args[i], Res: funType}
		}

		if err := unify(self.t, funType); err != nil {
			return err
		}

		// generalize, but never over the variables the environment
		// locks — including the function's own old binding
		env[decl.Name] = generalize(without(env, decl.Name), self.t)

		return nil
	}

	return fmt.Errorf("type error: unknown declaration")
}

// extend copies env with the extra bindings layered on top.
func extend(env map[string]*scheme, extra map[string]*scheme) map[string]*scheme {
	out := make(map[string]*scheme, len(env)+len(extra))

	for k, v := range env {
		out[k] = v
	}

	for k, v := range extra {
		out[k] = v
	}

	return out
}

// inferExpr types one expression.
func (c *typeChecker) inferExpr(env map[string]*scheme, e Expr) (Type, error) {
	switch expr := e.(type) {
	case IntLit:
		return intType, nil
	case BoolLit:
		return boolType, nil
	case Ident:
		if expr.Name == "~" {
			return &TFun{Arg: intType, Res: intType}, nil
		}

		s, bound := env[expr.Name]
		if !bound {
			return nil, fmt.Errorf("type error: unbound variable %s", expr.Name)
		}

		return c.instantiate(s), nil
	case If:
		cond, err := c.inferExpr(env, expr.Cond)
		if err != nil {
			return nil, err
		}

		if err := unify(cond, boolType); err != nil {
			return nil, err
		}

		then, err := c.inferExpr(env, expr.Then)
		if err != nil {
			return nil, err
		}

		els, err := c.inferExpr(env, expr.Else)
		if err != nil {
			return nil, err
		}

		if err := unify(then, els); err != nil {
			return nil, err
		}

		return then, nil
	case Fn:
		bindings := map[string]Type{}

		arg := c.patternType(expr.Arg, bindings)

		inner := map[string]*scheme{}

		for name, t := range bindings {
			inner[name] = &scheme{t: t}
		}

		body, err := c.inferExpr(extend(env, inner), expr.Body)
		if err != nil {
			return nil, err
		}

		return &TFun{Arg: arg, Res: body}, nil
	case App:
		fn, err := c.inferExpr(env, expr.Fn)
		if err != nil {
			return nil, err
		}

		arg, err := c.inferExpr(env, expr.Arg)
		if err != nil {
			return nil, err
		}

		res := c.freshVar()

		if err := unify(fn, &TFun{Arg: arg, Res: res}); err != nil {
			return nil, err
		}

		return res, nil
	case Infix:
		return c.inferInfix(env, expr)
	case Tuple:
		items := make([]Type, 0, len(expr.Items))

		for _, item := range expr.Items {
			t, err := c.inferExpr(env, item)
			if err != nil {
				return nil, err
			}

			items = append(items, t)
		}

		return &TTuple{Items: items}, nil
	case Let:
		inner := make(map[string]*scheme, len(env)+len(expr.Decls))

		for k, v := range env {
			inner[k] = v
		}

		for _, d := range expr.Decls {
			if err := c.inferDecl(inner, d, false); err != nil {
				return nil, err
			}
		}

		return c.inferExpr(inner, expr.Body)
	}

	return nil, fmt.Errorf("type error: unknown expression")
}

// inferInfix types an operator application: arithmetic on integers,
// comparisons of integers producing booleans, and the two boolean
// connectives.
func (c *typeChecker) inferInfix(env map[string]*scheme, e Infix) (Type, error) {
	l, err := c.inferExpr(env, e.L)
	if err != nil {
		return nil, err
	}

	r, err := c.inferExpr(env, e.R)
	if err != nil {
		return nil, err
	}

	switch e.Op {
	case "+", "-", "*", "/":
		if err := unify(l, intType); err != nil {
			return nil, err
		}

		if err := unify(r, intType); err != nil {
			return nil, err
		}

		return intType, nil
	case "=", "<>", "<", "<=", ">", ">=":
		if err := unify(l, intType); err != nil {
			return nil, err
		}

		if err := unify(r, intType); err != nil {
			return nil, err
		}

		return boolType, nil
	case "andalso", "orelse":
		if err := unify(l, boolType); err != nil {
			return nil, err
		}

		if err := unify(r, boolType); err != nil {
			return nil, err
		}

		return boolType, nil
	}

	return nil, fmt.Errorf("type error: unknown operator %s", e.Op)
}

// patternType shapes the type a pattern matches and collects the
// types of the variables it binds.
func (c *typeChecker) patternType(p Pat, bindings map[string]Type) Type {
	switch pat := p.(type) {
	case WildcardPat:
		return c.freshVar()
	case VarPat:
		t := c.freshVar()

		bindings[pat.Name] = t

		return t
	case ConstPat:
		if pat.Text == "true" || pat.Text == "false" {
			return boolType
		}

		return intType
	case TuplePat:
		items := make([]Type, 0, len(pat.Items))

		for _, item := range pat.Items {
			items = append(items, c.patternType(item, bindings))
		}

		return &TTuple{Items: items}
	}

	return c.freshVar()
}
