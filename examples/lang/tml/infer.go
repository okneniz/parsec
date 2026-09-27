package tml

import (
	"fmt"
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

	// pre-bind every top-level function with a fresh variable, so
	// that a body may reference functions declared later; the shared
	// variable is what carries the mutual constraints
	pos := 1

	for _, d := range decls {
		if _, ok := d.(FunDecl); ok {
			c.frame[pos] = &scheme{t: c.freshVar()}
			pos++

			continue
		}

		if val, ok := d.(ValDecl); ok {
			for _, b := range val.Binds {
				pos += len(patVars(b.Pat))
			}
		}
	}

	pos = 1

	for _, d := range decls {
		consumed, err := c.inferTopDecl(pos, d)
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

// inferTopDecl types one top-level declaration into the frame slots
// at pos and returns the slots consumed. A function name is already
// pre-bound there; a val pattern fills its slots with generalized
// schemes as they are inferred.
func (c *typeChecker) inferTopDecl(pos int, d Decl) (int, error) {
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
		self := c.frame[pos]

		// the parameters live in an inner scope: they must not
		// leak into the frame the declaration extends
		body := c.frame[:pos+1]

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

		// generalize, but never over the variables the environment
		// locks — the frame below the function's own slot
		c.frame[pos] = generalize(c.frame[:pos], self.t)

		return 1, nil
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

	for i := len(args) - 1; i >= 0; i-- {
		out = &TFun{Arg: args[i], Res: out}
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
	case ForwardRef:
		return nil, fmt.Errorf("type error: unbound variable %s", expr.Name)
	case TopVar:
		return c.instantiate(c.frame[expr.Pos]), nil
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
