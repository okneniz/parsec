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
// typechecker and the interpreter.
func evalProgram(decls []Decl) (Env, error) {
	if _, err := inferProgram(decls); err != nil {
		return nil, err
	}

	// one shared environment: closures capture it by reference, so
	// top-level functions see each other regardless of order
	env := Env{}

	for _, d := range decls {
		if err := evalDecl(env, d); err != nil {
			return nil, err
		}
	}

	return env, nil
}

// evalDecl evaluates one declaration into env.
func evalDecl(env Env, d Decl) error {
	switch decl := d.(type) {
	case ValDecl:
		for _, b := range decl.Binds {
			v, err := evalExpr(env, b.E)
			if err != nil {
				return err
			}

			bindings := Env{}

			if err := match(b.Pat, v, bindings); err != nil {
				return err
			}

			for name, value := range bindings {
				env[name] = value
			}
		}

		return nil
	case FunDecl:
		// the closure captures env, which will contain the function
		// itself by the time it is called
		env[decl.Name] = curry(decl.Args, decl.Body, env)

		return nil
	}

	return fmt.Errorf("runtime error: unknown declaration")
}

// curry wraps a multi-argument function body into nested one-argument
// closures.
func curry(args []Pat, body Expr, env Env) Value {
	if len(args) == 0 {
		return &VClosure{Arg: WildcardPat{}, Body: body, Env: env}
	}

	if len(args) == 1 {
		return &VClosure{Arg: args[0], Body: body, Env: env}
	}

	return &VClosure{
		Arg: args[0],
		Body: Fn{
			Arg:  args[1],
			Body: nestFn(args[2:], body),
		},
		Env: env,
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

// evalExpr evaluates one expression.
func evalExpr(env Env, e Expr) (Value, error) {
	switch expr := e.(type) {
	case IntLit:
		return VInt(expr.Value), nil
	case BoolLit:
		return VBool(expr.Value), nil
	case Ident:
		if v, bound := env[expr.Name]; bound {
			return v, nil
		}

		return nil, fmt.Errorf("runtime error: unbound variable %s", expr.Name)
	case If:
		cond, err := evalExpr(env, expr.Cond)
		if err != nil {
			return nil, err
		}

		if bool(cond.(VBool)) {
			return evalExpr(env, expr.Then)
		}

		return evalExpr(env, expr.Else)
	case Fn:
		return &VClosure{Arg: expr.Arg, Body: expr.Body, Env: env}, nil
	case App:
		// the lexer glued negation into App(~, literal)
		if id, ok := expr.Fn.(Ident); ok && id.Name == "~" {
			arg, err := evalExpr(env, expr.Arg)
			if err != nil {
				return nil, err
			}

			return VInt(-int64(arg.(VInt))), nil
		}

		fn, err := evalExpr(env, expr.Fn)
		if err != nil {
			return nil, err
		}

		arg, err := evalExpr(env, expr.Arg)
		if err != nil {
			return nil, err
		}

		return apply(fn, arg)
	case Infix:
		return evalInfix(env, expr)
	case Tuple:
		items := make(VTuple, 0, len(expr.Items))

		for _, item := range expr.Items {
			v, err := evalExpr(env, item)
			if err != nil {
				return nil, err
			}

			items = append(items, v)
		}

		return items, nil
	case Let:
		inner := make(Env, len(env)+len(expr.Decls))

		for k, v := range env {
			inner[k] = v
		}

		for _, d := range expr.Decls {
			if err := evalDecl(inner, d); err != nil {
				return nil, err
			}
		}

		return evalExpr(inner, expr.Body)
	}

	return nil, fmt.Errorf("runtime error: unknown expression")
}

// apply calls a closure: the argument is matched against the pattern
// and the body runs in the captured environment extended with the
// bindings.
func apply(fn Value, arg Value) (Value, error) {
	closure, ok := fn.(*VClosure)
	if !ok {
		return nil, fmt.Errorf("runtime error: not a function: %s", fn)
	}

	bindings := Env{}

	if err := match(closure.Arg, arg, bindings); err != nil {
		return nil, err
	}

	frame := make(Env, len(closure.Env)+len(bindings))

	for k, v := range closure.Env {
		frame[k] = v
	}

	for k, v := range bindings {
		frame[k] = v
	}

	return evalExpr(frame, closure.Body)
}

// evalInfix evaluates an operator application. The typechecker has
// already fixed the operand types, so the runtime only guards what it
// cannot know: the zero divisor.
func evalInfix(env Env, e Infix) (Value, error) {
	switch e.Op {
	case "andalso":
		l, err := evalExpr(env, e.L)
		if err != nil {
			return nil, err
		}

		if !bool(l.(VBool)) {
			return VBool(false), nil
		}

		r, err := evalExpr(env, e.R)
		if err != nil {
			return nil, err
		}

		return r, nil
	case "orelse":
		l, err := evalExpr(env, e.L)
		if err != nil {
			return nil, err
		}

		if bool(l.(VBool)) {
			return VBool(true), nil
		}

		return evalExpr(env, e.R)
	}

	l, err := evalExpr(env, e.L)
	if err != nil {
		return nil, err
	}

	r, err := evalExpr(env, e.R)
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

// match binds the pattern variables of p to the parts of v.
func match(p Pat, v Value, bindings Env) error {
	switch pat := p.(type) {
	case WildcardPat:
		return nil
	case VarPat:
		bindings[pat.Name] = v

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
			if err := match(item, tuple[i], bindings); err != nil {
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
