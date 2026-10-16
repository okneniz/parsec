package tml

import (
	"fmt"
	"os"
	"strconv"
)

// Eval parses, typechecks and interprets src and then calls its main
// function, returning the value main computes.
func Eval(code string) (Value, error) {
	ast, err := Parse(code)
	if err != nil {
		return nil, err
	}

	return evalProgram(ast)
}

// EvalFile reads a tiny ml program from a file and evaluates it.
func EvalFile(path string) (Value, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Eval(string(src))
}

// evalProgram runs the already parsed and already typechecked
// declarations and calls main. The functions become one shared row of
// knots — the builtin ~ at the bottom, every fun above it, each
// holding itself in its shadow slot — and the vals run in declaration
// order, each substituting its values into every function body and
// every later binding. The terms end up carrying no val slots: a free
// index of a running term always names a knot.
func evalProgram(decls []Decl) (Value, error) {
	if _, err := inferProgram(decls); err != nil {
		return nil, err
	}

	funs := make([]*VRec, 0)
	funAt := make([]int, 0)
	main := -1

	for i, d := range decls {
		decl, ok := d.(FunDecl)
		if !ok {
			continue
		}

		if decl.Name == "main" {
			main = len(funs)
		}

		funs = append(funs, &VRec{Args: decl.Args, Body: decl.Body})
		funAt = append(funAt, i)
	}

	row := make([]Value, 0, len(funs)+1)
	row = append(row, negBuiltin{})

	for _, v := range funs {
		row = append(row, v)
	}

	for _, v := range funs {
		v.Row = appendRow(row, v)
		v.Depth = len(v.Row)
	}

	site := len(row)

	for i, d := range decls {
		decl, ok := d.(ValDecl)
		if !ok {
			continue
		}

		for _, b := range decl.Binds {
			value, err := evalExpr(b.E, row, site)
			if err != nil {
				return nil, err
			}

			var values []Value

			if err := match(b.Pat, value, &values); err != nil {
				return nil, err
			}

			// the values leave the terms: everything still to run
			// reads them through the substitution instead of a slot.
			// A fun declared before the val never mentions it — its
			// body was parsed with the val not yet bound.
			for j, v := range funs {
				if funAt[j] < i {
					continue
				}

				substituteIntoFun(v, values, valsBetween(decls, i, funAt[j]))
			}

			for j := i + 1; j < len(decls); j++ {
				if val, isVal := decls[j].(ValDecl); isVal {
					decls[j] = substitutionTopsDecl(len(values), values, valsBetween(decls, i, j), val)
				}
			}
		}
	}

	if main < 0 {
		return nil, fmt.Errorf("runtime error: no main function")
	}

	return apply(funs[main], Tuple{})
}

// valsBetween counts the val binders the declarations strictly
// between from and to add. The top level hoists its funs below every
// val, so from the frame of a declaration at to only the vals sit
// above a val block declared at from.
func valsBetween(decls []Decl, from int, to int) int {
	out := 0

	for _, d := range decls[from+1 : to] {
		if val, ok := d.(ValDecl); ok {
			out += declArity(val)
		}
	}

	return out
}

// bindersBetween counts every binder the declarations strictly
// between from and to add. A let keeps declaration order — a fun
// name is a binder like any other — so all of them count.
func bindersBetween(decls []Decl, from int, to int) int {
	out := 0

	for _, d := range decls[from+1 : to] {
		out += declArity(d)
	}

	return out
}

// substituteIntoFun replaces the binders of one matched val in a
// function body: the body sits under the fun's own binders, which the
// declaration-level substitution crosses itself.
func substituteIntoFun(v *VRec, values []Value, under int) {
	decl := substitutionTopsDecl(len(values), values, under, FunDecl{
		Args: v.Args,
		Body: v.Body,
	}).(FunDecl)

	v.Body = decl.Body
}

// apply calls a function value: the argument is matched against the
// innermost argument pattern and its values are substituted into the
// body — the beta step — before the body runs against the row the
// function captured. The first argument is not the innermost binder
// of the body: the rest of the argument list sits above it, so the
// substitution enters at that depth. A multi-argument function left
// with arguments is again a function value.
func apply(fn Value, arg Value) (Value, error) {
	callee, ok := fn.(*VRec)
	if !ok || len(callee.Args) == 0 {
		return nil, fmt.Errorf("runtime error: not a function: %s", renderValue(fn))
	}

	var values []Value

	if err := match(callee.Args[0], arg, &values); err != nil {
		return nil, err
	}

	under := argsArity(callee.Args[1:])
	body := shift(-len(values), under, substTops(len(values), values, under, callee.Body))

	if len(callee.Args) > 1 {
		return &VRec{
			Args:  callee.Args[1:],
			Body:  body,
			Row:   callee.Row,
			Depth: callee.Depth,
		}, nil
	}

	return evalExpr(body, callee.Row, callee.Depth)
}

// evalExpr evaluates one expression against the knot row. The row is
// positional: its length is the depth the free indices of e count
// from, because every val binder has already left the terms.
func evalExpr(e Expr, row []Value, depth int) (Value, error) {
	switch expr := e.(type) {
	case IntLit:
		return expr, nil
	case BoolLit:
		return expr, nil
	case *VRec:
		return expr, nil
	case negBuiltin:
		return expr, nil
	case Var:
		pos := depth - 1 - expr.Index

		if pos < 0 || pos >= len(row) {
			return nil, fmt.Errorf("runtime error: unbound variable at index %d", expr.Index)
		}

		return row[pos], nil
	case FreeVar:
		return nil, fmt.Errorf("runtime error: unbound variable %s", expr.Name)
	case If:
		cond, err := evalExpr(expr.Cond, row, depth)
		if err != nil {
			return nil, err
		}

		if cond.(BoolLit).Value {
			return evalExpr(expr.Then, row, depth)
		}

		return evalExpr(expr.Else, row, depth)
	case Fn:
		// a lambda is a knot without a slot of its own: the body runs
		// against the row of the site, its argument substituted away
		return &VRec{
			Args:  []Pat{expr.Arg},
			Body:  expr.Body,
			Row:   row,
			Depth: depth,
		}, nil
	case App:
		fn, err := evalExpr(expr.Fn, row, depth)
		if err != nil {
			return nil, err
		}

		// the lexer glued negation into App(~, literal): ~ is the
		// builtin marker at the bottom of the row
		if _, neg := fn.(negBuiltin); neg {
			arg, err := evalExpr(expr.Arg, row, depth)
			if err != nil {
				return nil, err
			}

			n := arg.(IntLit)

			return IntLit{
				Text:  strconv.FormatInt(-n.Value, 10),
				Value: -n.Value,
			}, nil
		}

		arg, err := evalExpr(expr.Arg, row, depth)
		if err != nil {
			return nil, err
		}

		return apply(fn, arg)
	case Infix:
		return evalInfix(row, depth, expr)
	case Tuple:
		items := make([]Expr, 0, len(expr.Items))

		for _, item := range expr.Items {
			v, err := evalExpr(item, row, depth)
			if err != nil {
				return nil, err
			}

			items = append(items, v)
		}

		return Tuple{Items: items}, nil
	case Let:
		return evalLet(row, depth, expr)
	}

	return nil, fmt.Errorf("runtime error: unknown expression")
}

// evalLet runs the declarations of a let. The vals substitute their
// values into everything the rest of the block reads; every fun
// extends the row with its own knot and passes it on.
func evalLet(row []Value, depth int, expr Let) (Value, error) {
	innerRow, innerDepth := row, depth
	decls, body := expr.Decls, expr.Body
	funs := []*VRec{}
	funAt := []int{}

	for i, d := range decls {
		switch decl := d.(type) {
		case ValDecl:
			for _, b := range decl.Binds {
				value, err := evalExpr(b.E, innerRow, innerDepth)
				if err != nil {
					return nil, err
				}

				var values []Value

				if err := match(b.Pat, value, &values); err != nil {
					return nil, err
				}

				k := len(values)

				for j, v := range funs {
					if funAt[j] < i {
						continue
					}

					substituteIntoFun(v, values, bindersBetween(decls, i, funAt[j]))
				}

				for j := i + 1; j < len(decls); j++ {
					decls[j] = substitutionTopsDecl(k, values, bindersBetween(decls, i, j), decls[j])
				}

				under := bindersBetween(decls, i, len(decls))
				body = shift(-k, under, substTops(k, values, under, body))
			}
		case FunDecl:
			v := &VRec{Args: decl.Args, Body: decl.Body}

			v.Row = appendRow(innerRow, v)
			v.Depth = len(v.Row)

			innerRow = appendRow(innerRow, v)
			innerDepth = len(innerRow)

			funs = append(funs, v)
			funAt = append(funAt, i)
		}
	}

	return evalExpr(body, innerRow, innerDepth)
}

// evalInfix evaluates an operator application. The typechecker has
// already fixed the operand types, so the runtime only guards what it
// cannot know: the zero divisor.
func evalInfix(row []Value, depth int, e Infix) (Value, error) {
	switch e.Op {
	case "andalso":
		l, err := evalExpr(e.L, row, depth)
		if err != nil {
			return nil, err
		}

		if !l.(BoolLit).Value {
			return BoolLit{Value: false}, nil
		}

		r, err := evalExpr(e.R, row, depth)
		if err != nil {
			return nil, err
		}

		return r, nil
	case "orelse":
		l, err := evalExpr(e.L, row, depth)
		if err != nil {
			return nil, err
		}

		if l.(BoolLit).Value {
			return BoolLit{Value: true}, nil
		}

		return evalExpr(e.R, row, depth)
	}

	l, err := evalExpr(e.L, row, depth)
	if err != nil {
		return nil, err
	}

	r, err := evalExpr(e.R, row, depth)
	if err != nil {
		return nil, err
	}

	li, ri := l.(IntLit).Value, r.(IntLit).Value

	switch e.Op {
	case "+":
		return IntLit{Text: strconv.FormatInt(li+ri, 10), Value: li + ri}, nil
	case "-":
		return IntLit{Text: strconv.FormatInt(li-ri, 10), Value: li - ri}, nil
	case "*":
		return IntLit{Text: strconv.FormatInt(li*ri, 10), Value: li * ri}, nil
	case "/":
		if ri == 0 {
			return nil, fmt.Errorf("runtime error: division by zero")
		}

		return IntLit{Text: strconv.FormatInt(li/ri, 10), Value: li / ri}, nil
	case "=":
		return BoolLit{Value: li == ri}, nil
	case "<>":
		return BoolLit{Value: li != ri}, nil
	case "<":
		return BoolLit{Value: li < ri}, nil
	case "<=":
		return BoolLit{Value: li <= ri}, nil
	case ">":
		return BoolLit{Value: li > ri}, nil
	case ">=":
		return BoolLit{Value: li >= ri}, nil
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
		expected := Value(IntLit{Text: pat.Text, Value: parseInt64(pat.Text)})

		switch pat.Text {
		case "true":
			expected = BoolLit{Value: true}
		case "false":
			expected = BoolLit{Value: false}
		}

		equal := false

		switch got := v.(type) {
		case IntLit:
			want, isInt := expected.(IntLit)
			equal = isInt && got == want
		case BoolLit:
			want, isBool := expected.(BoolLit)
			equal = isBool && got == want
		}

		if !equal {
			return fmt.Errorf("runtime error: pattern match failure: %s vs %s", p, renderValue(v))
		}

		return nil
	case TuplePat:
		tuple, ok := v.(Tuple)
		if !ok || len(tuple.Items) != len(pat.Items) {
			return fmt.Errorf("runtime error: pattern match failure: %s vs %s", p, renderValue(v))
		}

		for i, item := range pat.Items {
			if err := match(item, tuple.Items[i], values); err != nil {
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
