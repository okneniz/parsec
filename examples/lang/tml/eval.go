package tml

import (
	"fmt"
	"os"
)

// Eval parses, typechecks and interprets src, returning the values of
// the top-level bindings.
func Eval(src string) (Env, error) {
	decls, err := Parse(src)
	if err != nil {
		return nil, err
	}

	return evalProgram(decls)
}

// EvalFile reads a tiny ml program from a file and evaluates it.
func EvalFile(path string) (Env, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Eval(string(src))
}

// evalProgram runs the already parsed and already typechecked
// declarations — Eval parses once and passes the result to both the
// typechecker and the interpreter. The environment mirrors the
// binder stack the parser kept: one cell per binder, and the
// top-level cells live in a shared frame, patched as the
// declarations run.
func evalProgram(decls []Decl) (Env, error) {
	if _, err := inferProgram(decls); err != nil {
		return nil, err
	}

	names := layout(decls)

	frame := make([]*cell, len(names))
	for i := range frame {
		frame[i] = &cell{}
	}

	pos := 1

	for _, d := range decls {
		consumed, err := evalTopDecl(frame, pos, d)
		if err != nil {
			return nil, err
		}

		pos += consumed
	}

	env := Env{}

	for i, name := range names {
		if name == "~" {
			continue
		}

		env[name] = frame[i].val
	}

	return env, nil
}

// evalTopDecl evaluates one top-level declaration into the frame
// slots at pos and returns the slots consumed.
func evalTopDecl(frame []*cell, pos int, d Decl) (int, error) {
	switch decl := d.(type) {
	case ValDecl:
		consumed := 0

		for _, b := range decl.Binds {
			v, err := evalExpr(frame, frame[:pos], b.E)
			if err != nil {
				return 0, err
			}

			var values []Value

			if err := match(b.Pat, v, &values); err != nil {
				return 0, err
			}

			for _, value := range values {
				frame[pos].val = value
				pos++
				consumed++
			}
		}

		return consumed, nil
	case FunDecl:
		// the closure reads the frame through its own slot, which
		// holds the closure itself by the time it is called
		frame[pos].val = curry(frame, frame[:pos+1], decl.Args, decl.Body)

		return 1, nil
	}

	return 0, fmt.Errorf("runtime error: unknown declaration")
}

// evalLocalDecl evaluates one declaration of a let and returns the
// local chain extended with the bound cells.
func evalLocalDecl(frame []*cell, env []*cell, d Decl) ([]*cell, error) {
	switch decl := d.(type) {
	case ValDecl:
		for _, b := range decl.Binds {
			v, err := evalExpr(frame, env, b.E)
			if err != nil {
				return nil, err
			}

			var values []Value

			if err := match(b.Pat, v, &values); err != nil {
				return nil, err
			}

			for _, value := range values {
				env = appendCell(env, &cell{val: value})
			}
		}

		return env, nil
	case FunDecl:
		// the cell comes first, the closure fills it after it is
		// built: the recursion of a let-bound function
		self := &cell{}

		env = appendCell(env, self)
		self.val = curry(frame, env, decl.Args, decl.Body)

		return env, nil
	}

	return nil, fmt.Errorf("runtime error: unknown declaration")
}

// appendCell copies env with one more cell on top: the snapshots
// never share a backing array, so the scopes stay their own.
func appendCell(env []*cell, c *cell) []*cell {
	out := make([]*cell, len(env), len(env)+1)
	copy(out, env)

	return append(out, c)
}

// curry wraps a multi-argument function body into nested one-argument
// closures.
func curry(frame []*cell, env []*cell, args []Pat, body Expr) Value {
	if len(args) == 0 {
		return &VClosure{Arg: WildcardPat{}, Body: body, Frame: frame, Env: env}
	}

	if len(args) == 1 {
		return &VClosure{Arg: args[0], Body: body, Frame: frame, Env: env}
	}

	return &VClosure{
		Arg: args[0],
		Body: Fn{
			Arg:  args[1],
			Body: nestFn(args[2:], body),
		},
		Frame: frame,
		Env:   env,
	}
}

// nestFn builds the nested fn-expression for the tail of a curried
// argument list.
func nestFn(args []Pat, body Expr) Expr {
	if len(args) == 0 {
		return body
	}

	return Fn{Arg: args[0], Body: nestFn(args[1:], body)}
}

// evalExpr evaluates one expression. The environment is the binder
// stack the parser had: one cell per binder, innermost last.
func evalExpr(frame []*cell, env []*cell, e Expr) (Value, error) {
	switch expr := e.(type) {
	case IntLit:
		return VInt(expr.Value), nil
	case BoolLit:
		return VBool(expr.Value), nil
	case Var:
		if v := env[len(env)-1-expr.Index].val; v != nil {
			return v, nil
		}

		return nil, fmt.Errorf("runtime error: unbound variable at index %d", expr.Index)
	case TopVar:
		if v := frame[expr.Pos].val; v != nil {
			return v, nil
		}

		return nil, fmt.Errorf("runtime error: unbound variable at top-level slot %d", expr.Pos)
	case If:
		cond, err := evalExpr(frame, env, expr.Cond)
		if err != nil {
			return nil, err
		}

		if bool(cond.(VBool)) {
			return evalExpr(frame, env, expr.Then)
		}

		return evalExpr(frame, env, expr.Else)
	case Fn:
		return &VClosure{Arg: expr.Arg, Body: expr.Body, Frame: frame, Env: env}, nil
	case App:
		// the lexer glued negation into App(~, literal): ~ is the
		// builtin binder at the bottom of the environment
		if v, ok := expr.Fn.(Var); ok && v.Index == len(env)-1 {
			arg, err := evalExpr(frame, env, expr.Arg)
			if err != nil {
				return nil, err
			}

			return VInt(-int64(arg.(VInt))), nil
		}

		fn, err := evalExpr(frame, env, expr.Fn)
		if err != nil {
			return nil, err
		}

		arg, err := evalExpr(frame, env, expr.Arg)
		if err != nil {
			return nil, err
		}

		return apply(frame, fn, arg)
	case Infix:
		return evalInfix(frame, env, expr)
	case Tuple:
		items := make(VTuple, 0, len(expr.Items))

		for _, item := range expr.Items {
			v, err := evalExpr(frame, env, item)
			if err != nil {
				return nil, err
			}

			items = append(items, v)
		}

		return items, nil
	case Let:
		inner := env

		for _, d := range expr.Decls {
			extended, err := evalLocalDecl(frame, inner, d)
			if err != nil {
				return nil, err
			}

			inner = extended
		}

		return evalExpr(frame, inner, expr.Body)
	}

	return nil, fmt.Errorf("runtime error: unknown expression")
}

// apply calls a closure: the argument is matched against the pattern
// and the body runs in the captured environment extended with the
// bindings.
func apply(frame []*cell, fn Value, arg Value) (Value, error) {
	closure, ok := fn.(*VClosure)
	if !ok {
		return nil, fmt.Errorf("runtime error: not a function: %s", fn)
	}

	var values []Value

	if err := match(closure.Arg, arg, &values); err != nil {
		return nil, err
	}

	env := make([]*cell, len(closure.Env), len(closure.Env)+len(values))
	copy(env, closure.Env)

	for _, value := range values {
		env = append(env, &cell{val: value})
	}

	return evalExpr(closure.Frame, env, closure.Body)
}

// evalInfix evaluates an operator application. The typechecker has
// already fixed the operand types, so the runtime only guards what it
// cannot know: the zero divisor.
func evalInfix(frame []*cell, env []*cell, e Infix) (Value, error) {
	switch e.Op {
	case "andalso":
		l, err := evalExpr(frame, env, e.L)
		if err != nil {
			return nil, err
		}

		if !bool(l.(VBool)) {
			return VBool(false), nil
		}

		r, err := evalExpr(frame, env, e.R)
		if err != nil {
			return nil, err
		}

		return r, nil
	case "orelse":
		l, err := evalExpr(frame, env, e.L)
		if err != nil {
			return nil, err
		}

		if bool(l.(VBool)) {
			return VBool(true), nil
		}

		return evalExpr(frame, env, e.R)
	}

	l, err := evalExpr(frame, env, e.L)
	if err != nil {
		return nil, err
	}

	r, err := evalExpr(frame, env, e.R)
	if err != nil {
		return nil, err
	}

	li, ri := int64(l.(VInt)), int64(r.(VInt))

	switch e.Op {
	case "+":
		return VInt(li + ri), nil
	case "-":
		return VInt(li - ri), nil
	case "*":
		return VInt(li * ri), nil
	case "/":
		if ri == 0 {
			return nil, fmt.Errorf("runtime error: division by zero")
		}

		return VInt(li / ri), nil
	case "=":
		return VBool(li == ri), nil
	case "<>":
		return VBool(li != ri), nil
	case "<":
		return VBool(li < ri), nil
	case "<=":
		return VBool(li <= ri), nil
	case ">":
		return VBool(li > ri), nil
	case ">=":
		return VBool(li >= ri), nil
	}

	return nil, fmt.Errorf("runtime error: unknown operator %s", e.Op)
}

// match binds the pattern variables of p to the parts of v, in
// binding order.
func match(p Pat, v Value, values *[]Value) error {
	switch pat := p.(type) {
	case WildcardPat:
		return nil
	case VarPat:
		*values = append(*values, v)

		return nil
	case ConstPat:
		expected := Value(VInt(parseInt64(pat.Text)))

		switch pat.Text {
		case "true":
			expected = VBool(true)
		case "false":
			expected = VBool(false)
		}

		if v != expected {
			return fmt.Errorf("runtime error: pattern match failure: %s vs %s", p, v)
		}

		return nil
	case TuplePat:
		tuple, ok := v.(VTuple)
		if !ok || len(tuple) != len(pat.Items) {
			return fmt.Errorf("runtime error: pattern match failure: %s vs %s", p, v)
		}

		for i, item := range pat.Items {
			if err := match(item, tuple[i], values); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("runtime error: unknown pattern")
}

// parseInt64 parses the text of an integer literal, with the ~
// negation of negative patterns resolved.
func parseInt64(text string) int64 {
	negated := false

	for len(text) > 0 && text[0] == '~' {
		negated = !negated
		text = text[1:]
	}

	var out int64

	for _, r := range text {
		out = out*10 + int64(r-'0')
	}

	if negated {
		return -out
	}

	return out
}
