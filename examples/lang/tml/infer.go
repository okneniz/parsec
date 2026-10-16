package tml

import (
	"fmt"
	"slices"
)

// typeChecker is one run of Hindley-Milner inference: a counter for
// fresh type variables and the top-level frame of schemes. The
// substitution lives in the variables themselves.
type typeChecker struct {
	fresh int
	frame []*scheme
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

// inferProgram typechecks already parsed declarations. The type
// environment mirrors the binder stack the parser kept: one scheme
// per binder, innermost last, over the top-level frame the program
// addresses positionally.
func inferProgram(decls []Decl) (map[string]Type, error) {
	c := &typeChecker{}
	names := layout(decls)

	c.frame = make([]*scheme, len(names))

	// the builtin ~ of the initial stack: int -> int
	c.frame[0] = &scheme{t: &TFun{Arg: intType, Res: intType}}

	// pre-bind every top-level function with a fresh variable, in the
	// fun row the second parse pass laid out: the shared variable is
	// what carries the mutual constraints, because every body sees
	// the whole row
	funs := 0

	for _, d := range decls {
		if _, ok := d.(FunDecl); ok {
			funs++
			c.frame[funs] = &scheme{t: c.freshVar()}
		}
	}

	// the slot every declaration starts from: the fun row is laid out
	// already, the val binders pile up after it in declaration order
	pos := funs + 1
	fun := 0

	for _, d := range decls {
		slot := 0

		if _, ok := d.(FunDecl); ok {
			fun++
			slot = fun
		}

		consumed, err := c.inferTopDecl(pos, slot, d)
		if err != nil {
			return nil, err
		}

		pos += consumed
	}

	out := map[string]Type{}

	for i, name := range names {
		if name == "~" {
			continue
		}

		out[name] = prune(c.frame[i].t)
	}

	return out, nil
}

// layout lists the binders of the top level in stack order: the
// builtin first, then every function — the pre-bound row of the
// second parse pass — then the val binders in declaration order. The
// typechecker addresses the row positionally; the reporting map of
// Infer reads its names.
func layout(decls []Decl) []string {
	names := []string{"~"}

	for _, d := range decls {
		if fun, ok := d.(FunDecl); ok {
			names = append(names, fun.Name)
		}
	}

	for _, d := range decls {
		if val, ok := d.(ValDecl); ok {
			for _, b := range val.Binds {
				names = append(names, patVars(b.Pat)...)
			}
		}
	}

	return names
}

// inferTopDecl types one top-level declaration into its frame slots —
// pos for a val, slot for a fun — and returns the slots consumed. A
// function name is already pre-bound in the fun row; a val pattern
// fills its slots with generalized schemes as they are inferred.
func (c *typeChecker) inferTopDecl(pos int, slot int, d Decl) (int, error) {
	switch decl := d.(type) {
	case ValDecl:
		consumed := 0

		for _, b := range decl.Binds {
			outer := c.frame[:pos]

			patType, bindings := c.patternType(b.Pat)

			valueType, err := c.inferExpr(outer, b.E)
			if err != nil {
				return 0, err
			}

			if err := unify(patType, valueType); err != nil {
				return 0, err
			}

			// let-polymorphism: each variable of the pattern is
			// generalized over what the environment does not lock
			for _, t := range bindings {
				c.frame[pos] = generalize(outer, t)
				pos++
				consumed++
			}
		}

		return consumed, nil
	case FunDecl:
		self := c.frame[slot]

		// the body sees the pre-bound slot again, as its own shadow
		// innermost of the base, with the arguments above it; the
		// parameters live in an inner scope and must not leak into
		// the frame the declaration extends
		body := appendEnv(c.frame[:pos], self)

		var argTypes []Type

		for _, arg := range decl.Args {
			t, bindings := c.patternType(arg)

			argTypes = append(argTypes, t)

			for _, bt := range bindings {
				body = appendEnv(body, &scheme{t: bt})
			}
		}

		bodyType, err := c.inferExpr(body, decl.Body)
		if err != nil {
			return 0, err
		}

		funType := funTypeOf(argTypes, bodyType)

		if err := unify(self.t, funType); err != nil {
			return 0, err
		}

		// generalize over the base the body saw — the fun row and the
		// vals before it — minus the own slot: its raw variable would
		// lock the whole type and make every recursive fun monomorphic
		locked := make([]*scheme, 0, pos-1)
		locked = append(locked, c.frame[:slot]...)
		locked = append(locked, c.frame[slot+1:pos]...)

		c.frame[slot] = generalize(locked, self.t)

		return 0, nil
	}

	return 0, fmt.Errorf("type error: unknown declaration")
}

// inferLocalDecl types one declaration of a let and returns the
// local chain extended with the bound schemes.
func (c *typeChecker) inferLocalDecl(env []*scheme, d Decl) ([]*scheme, error) {
	switch decl := d.(type) {
	case ValDecl:
		for _, b := range decl.Binds {
			patType, bindings := c.patternType(b.Pat)

			valueType, err := c.inferExpr(env, b.E)
			if err != nil {
				return nil, err
			}

			if err := unify(patType, valueType); err != nil {
				return nil, err
			}

			for _, t := range bindings {
				env = appendEnv(env, generalize(env, t))
			}
		}

		return env, nil
	case FunDecl:
		// the function name is a fresh scheme under the let,
		// shadowing any outer binding; its slot is written after
		// the body has fixed the type. The parameters scope over
		// the body only — the chain the let continues from ends
		// at the name
		outer := env

		env = appendEnv(env, &scheme{t: c.freshVar()})

		self := len(env) - 1
		afterSelf := env

		var argTypes []Type

		for _, arg := range decl.Args {
			t, bindings := c.patternType(arg)

			argTypes = append(argTypes, t)

			for _, bt := range bindings {
				env = appendEnv(env, &scheme{t: bt})
			}
		}

		bodyType, err := c.inferExpr(env, decl.Body)
		if err != nil {
			return nil, err
		}

		funType := funTypeOf(argTypes, bodyType)

		if err := unify(afterSelf[self].t, funType); err != nil {
			return nil, err
		}

		afterSelf[self] = generalize(outer, afterSelf[self].t)

		return afterSelf, nil
	}

	return nil, fmt.Errorf("type error: unknown declaration")
}

// funTypeOf folds the argument types over the body type.
func funTypeOf(args []Type, body Type) Type {
	out := body

	for _, arg := range slices.Backward(args) {
		out = &TFun{Arg: arg, Res: out}
	}

	return out
}

// appendEnv copies env with one more scheme on top: the snapshots
// never share a backing array, so the scopes stay their own.
func appendEnv(env []*scheme, s *scheme) []*scheme {
	out := make([]*scheme, len(env), len(env)+1)
	copy(out, env)

	return append(out, s)
}

// inferExpr types one expression. The environment is the binder
// stack the parser had: a scheme per binder, innermost last.
func (c *typeChecker) inferExpr(env []*scheme, e Expr) (Type, error) {
	switch expr := e.(type) {
	case IntLit:
		return intType, nil
	case BoolLit:
		return boolType, nil
	case Var:
		return c.instantiate(env[len(env)-1-expr.Index]), nil
	case FreeVar:
		return nil, fmt.Errorf("type error: unbound variable %s", expr.Name)
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
		arg, bindings := c.patternType(expr.Arg)

		inner := env

		for _, t := range bindings {
			inner = appendEnv(inner, &scheme{t: t})
		}

		body, err := c.inferExpr(inner, expr.Body)
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
		inner := env

		for _, d := range expr.Decls {
			extended, err := c.inferLocalDecl(inner, d)
			if err != nil {
				return nil, err
			}

			inner = extended
		}

		return c.inferExpr(inner, expr.Body)
	}

	return nil, fmt.Errorf("type error: unknown expression")
}

// inferInfix types an operator application: arithmetic on integers,
// comparisons of integers producing booleans, and the two boolean
// connectives.
func (c *typeChecker) inferInfix(env []*scheme, e Infix) (Type, error) {
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

// patternType shapes the type a pattern matches and lists the types
// of the variables it binds, in binding order.
func (c *typeChecker) patternType(p Pat) (Type, []Type) {
	switch pat := p.(type) {
	case WildcardPat:
		return c.freshVar(), nil
	case VarPat:
		t := c.freshVar()

		return t, []Type{t}
	case ConstPat:
		if pat.Text == "true" || pat.Text == "false" {
			return boolType, nil
		}

		return intType, nil
	case TuplePat:
		items := make([]Type, 0, len(pat.Items))

		var bindings []Type

		for _, item := range pat.Items {
			t, bs := c.patternType(item)

			items = append(items, t)
			bindings = append(bindings, bs...)
		}

		return &TTuple{Items: items}, bindings
	}

	return c.freshVar(), nil
}
