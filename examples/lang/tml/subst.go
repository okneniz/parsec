package tml

// The de Bruijn substitution toolkit of the interpreter: shift moves
// the free indices of a term across binders, subst replaces one index
// with a term, and substitutionTops substitutes the values a pattern
// matched for the binders of that pattern, removing them. Tiny ml
// patterns are first-order — they bind without reading — so none of
// this descends into a pattern: only its binder count matters.

// shift adds d to the index of every variable at or above the cutoff,
// leaving the variables the enclosing binders own alone.
func shift(d int, cutoff int, e Expr) Expr {
	switch expr := e.(type) {
	case Var:
		if expr.Index >= cutoff {
			return Var{Index: expr.Index + d}
		}

		return expr
	case If:
		return If{
			Cond: shift(d, cutoff, expr.Cond),
			Then: shift(d, cutoff, expr.Then),
			Else: shift(d, cutoff, expr.Else),
		}
	case Fn:
		return Fn{
			Arg:  expr.Arg,
			Body: shift(d, cutoff+patArity(expr.Arg), expr.Body),
		}
	case App:
		return App{
			Fn:  shift(d, cutoff, expr.Fn),
			Arg: shift(d, cutoff, expr.Arg),
		}
	case Infix:
		return Infix{
			Op: expr.Op,
			L:  shift(d, cutoff, expr.L),
			R:  shift(d, cutoff, expr.R),
		}
	case Tuple:
		items := make([]Expr, 0, len(expr.Items))

		for _, item := range expr.Items {
			items = append(items, shift(d, cutoff, item))
		}

		return Tuple{Items: items}
	case Let:
		return Let{
			Decls: shiftDecls(d, cutoff, expr.Decls),
			Body:  shift(d, cutoff+declsArity(expr.Decls), expr.Body),
		}
	}

	return e
}

// subst replaces every Var(j) with v, which the descent shifts past
// the binders it enters itself, and leaves the other indices for the
// caller to adjust.
func subst(j int, v Expr, e Expr) Expr {
	switch expr := e.(type) {
	case Var:
		if expr.Index == j {
			return v
		}

		return expr
	case If:
		return If{
			Cond: subst(j, v, expr.Cond),
			Then: subst(j, v, expr.Then),
			Else: subst(j, v, expr.Else),
		}
	case Fn:
		inner := patArity(expr.Arg)

		return Fn{
			Arg:  expr.Arg,
			Body: subst(j+inner, shift(inner, 0, v), expr.Body),
		}
	case App:
		return App{
			Fn:  subst(j, v, expr.Fn),
			Arg: subst(j, v, expr.Arg),
		}
	case Infix:
		return Infix{
			Op: expr.Op,
			L:  subst(j, v, expr.L),
			R:  subst(j, v, expr.R),
		}
	case Tuple:
		items := make([]Expr, 0, len(expr.Items))

		for _, item := range expr.Items {
			items = append(items, subst(j, v, item))
		}

		return Tuple{Items: items}
	case Let:
		return Let{
			Decls: substDecls(j, v, expr.Decls),
			Body:  subst(j+declsArity(expr.Decls), shift(declsArity(expr.Decls), 0, v), expr.Body),
		}
	}

	return e
}

// substitutionTop replaces the innermost binder of e with v and
// removes the binder: the classic beta step of a one-variable lambda.
func substitutionTop(e Expr, v Expr) Expr {
	return substitutionTops(e, []Expr{v})
}

// substitutionTops substitutes the values vs for the innermost
// len(vs) binders of e in one step — a pattern binds several — in
// binding order: the last value belongs to the innermost binder. The
// binders are removed, so the indices above them drop by len(vs).
func substitutionTops(e Expr, vs []Expr) Expr {
	return shift(-len(vs), 0, substTops(len(vs), vs, 0, e))
}

// substTops replaces the top k binders with vs simultaneously,
// leaving the indices above them for the caller to drop.
func substTops(k int, vs []Expr, cutoff int, e Expr) Expr {
	switch expr := e.(type) {
	case Var:
		if expr.Index < cutoff || expr.Index-cutoff >= k {
			return expr
		}

		// the value enters from outside the replaced block: it clears
		// the binders of the block and the ones above the site, and
		// the caller's shift drops the block again
		return shift(cutoff+k, 0, vs[k-1-(expr.Index-cutoff)])
	case If:
		return If{
			Cond: substTops(k, vs, cutoff, expr.Cond),
			Then: substTops(k, vs, cutoff, expr.Then),
			Else: substTops(k, vs, cutoff, expr.Else),
		}
	case Fn:
		return Fn{
			Arg:  expr.Arg,
			Body: substTops(k, vs, cutoff+patArity(expr.Arg), expr.Body),
		}
	case App:
		return App{
			Fn:  substTops(k, vs, cutoff, expr.Fn),
			Arg: substTops(k, vs, cutoff, expr.Arg),
		}
	case Infix:
		return Infix{
			Op: expr.Op,
			L:  substTops(k, vs, cutoff, expr.L),
			R:  substTops(k, vs, cutoff, expr.R),
		}
	case Tuple:
		items := make([]Expr, 0, len(expr.Items))

		for _, item := range expr.Items {
			items = append(items, substTops(k, vs, cutoff, item))
		}

		return Tuple{Items: items}
	case Let:
		return Let{
			Decls: substTopsDecls(k, vs, cutoff, expr.Decls),
			Body:  substTops(k, vs, cutoff+declsArity(expr.Decls), expr.Body),
		}
	}

	return e
}

// shiftDecls shifts the expressions of every declaration, each one
// under the binders of the declarations before it.
func shiftDecls(d int, cutoff int, decls []Decl) []Decl {
	out := make([]Decl, 0, len(decls))
	inner := cutoff

	for _, decl := range decls {
		out = append(out, shiftDecl(d, inner, decl))
		inner += declArity(decl)
	}

	return out
}

// substDecls substitutes through the declarations, raising j past the
// binders of each declaration as it goes.
func substDecls(j int, v Expr, decls []Decl) []Decl {
	out := make([]Decl, 0, len(decls))
	inner, moved := j, v

	for _, decl := range decls {
		out = append(out, substDecl(inner, moved, decl))

		arity := declArity(decl)
		inner += arity
		moved = shift(arity, 0, moved)
	}

	return out
}

// substTopsDecls substitutes the top binders through the declarations,
// raising the cutoff past the binders of each one as it goes.
func substTopsDecls(k int, vs []Expr, cutoff int, decls []Decl) []Decl {
	out := make([]Decl, 0, len(decls))
	inner := cutoff

	for _, decl := range decls {
		out = append(out, substTopsDecl(k, vs, inner, decl))
		inner += declArity(decl)
	}

	return out
}

// shiftDecl shifts the expression of one declaration; a fun binds its
// name for what follows and its arguments inside its own body only.
func shiftDecl(d int, cutoff int, decl Decl) Decl {
	switch dl := decl.(type) {
	case ValDecl:
		binds := make([]ValBind, 0, len(dl.Binds))
		inner := cutoff

		for _, b := range dl.Binds {
			binds = append(binds, ValBind{
				Pat: b.Pat,
				E:   shift(d, inner, b.E),
			})

			inner += patArity(b.Pat)
		}

		return ValDecl{Binds: binds}
	case FunDecl:
		return FunDecl{
			Name: dl.Name,
			Args: dl.Args,
			Body: shift(d, cutoff+1+argsArity(dl.Args), dl.Body),
		}
	}

	return decl
}

// substDecl substitutes through one declaration.
func substDecl(j int, v Expr, decl Decl) Decl {
	switch dl := decl.(type) {
	case ValDecl:
		binds := make([]ValBind, 0, len(dl.Binds))
		inner, moved := j, v

		for _, b := range dl.Binds {
			binds = append(binds, ValBind{
				Pat: b.Pat,
				E:   subst(inner, moved, b.E),
			})

			arity := patArity(b.Pat)
			inner += arity
			moved = shift(arity, 0, moved)
		}

		return ValDecl{Binds: binds}
	case FunDecl:
		under := 1 + argsArity(dl.Args)

		return FunDecl{
			Name: dl.Name,
			Args: dl.Args,
			Body: subst(j+under, shift(under, 0, v), dl.Body),
		}
	}

	return decl
}

// substTopsDecl substitutes the top binders through one declaration,
// leaving the removal of the replaced block to the caller.
func substTopsDecl(k int, vs []Expr, cutoff int, decl Decl) Decl {
	switch dl := decl.(type) {
	case ValDecl:
		binds := make([]ValBind, 0, len(dl.Binds))
		inner := cutoff

		for _, b := range dl.Binds {
			binds = append(binds, ValBind{
				Pat: b.Pat,
				E:   substTops(k, vs, inner, b.E),
			})

			inner += patArity(b.Pat)
		}

		return ValDecl{Binds: binds}
	case FunDecl:
		return FunDecl{
			Name: dl.Name,
			Args: dl.Args,
			Body: substTops(k, vs, cutoff+1+argsArity(dl.Args), dl.Body),
		}
	}

	return decl
}

// substitutionTopsDecl is substitutionTops through one declaration:
// the replaced block leaves every expression of the declaration, each
// at the depth that expression lives at.
func substitutionTopsDecl(k int, vs []Expr, under int, decl Decl) Decl {
	switch dl := decl.(type) {
	case ValDecl:
		binds := make([]ValBind, 0, len(dl.Binds))
		inner := under

		for _, b := range dl.Binds {
			binds = append(binds, ValBind{
				Pat: b.Pat,
				E:   shift(-k, inner, substTops(k, vs, inner, b.E)),
			})

			inner += patArity(b.Pat)
		}

		return ValDecl{Binds: binds}
	case FunDecl:
		inner := under + 1 + argsArity(dl.Args)

		return FunDecl{
			Name: dl.Name,
			Args: dl.Args,
			Body: shift(-k, inner, substTops(k, vs, inner, dl.Body)),
		}
	}

	return decl
}

// declsArity counts the binders a sequence of declarations adds.
func declsArity(decls []Decl) int {
	out := 0

	for _, decl := range decls {
		out += declArity(decl)
	}

	return out
}

// declArity counts the binders one declaration adds: every variable
// of every val pattern, or the single name of a fun.
func declArity(decl Decl) int {
	switch dl := decl.(type) {
	case ValDecl:
		out := 0

		for _, b := range dl.Binds {
			out += patArity(b.Pat)
		}

		return out
	case FunDecl:
		return 1
	}

	return 0
}

// argsArity counts the binders a function argument list adds.
func argsArity(args []Pat) int {
	out := 0

	for _, arg := range args {
		out += patArity(arg)
	}

	return out
}

// patArity counts the variables a pattern binds.
func patArity(p Pat) int {
	switch pat := p.(type) {
	case VarPat:
		return 1
	case TuplePat:
		out := 0

		for _, item := range pat.Items {
			out += patArity(item)
		}

		return out
	}

	return 0
}
