package tml

import "strconv"

// Expr is an expression node.
type Expr interface {
	String() string
}

// IntLit is an integer literal; Value is the parsed Text, best
// effort, the text itself is authoritative.
type IntLit struct {
	Text  string
	Value int64
}

// BoolLit is a boolean literal.
type BoolLit struct {
	Value bool
}

// Ident is a variable reference; tiny ml has no qualified names.
type Ident struct {
	Name string
}

// If is a conditional expression.
type If struct {
	Cond, Then, Else Expr
}

// Fn is a single-clause function expression: fn pat => body.
type Fn struct {
	Arg  Pat
	Body Expr
}

// App is function application; application binds tighter than any
// infix operator.
type App struct {
	Fn  Expr
	Arg Expr
}

// Infix is an infix application; Op is the operator lexeme, which may
// also be the reserved words andalso and orelse.
type Infix struct {
	Op   string
	L, R Expr
}

// Tuple is a tuple expression (e1, ..., en).
type Tuple struct {
	Items []Expr
}

// Let is a let expression with a single-expression body.
type Let struct {
	Decls []Decl
	Body  Expr
}

// Decl is a declaration node.
type Decl interface {
	String() string
}

// ValBind is one pat = exp binding of a val declaration.
type ValBind struct {
	Pat Pat
	E   Expr
}

// ValDecl is a val declaration.
type ValDecl struct {
	Binds []ValBind
}

// FunDecl is a recursive function declaration: fun f x1 ... xn = body,
// with atomic argument patterns.
type FunDecl struct {
	Name string
	Args []Pat
	Body Expr
}

// Pat is a pattern node. Tiny ml patterns are first-order: variables,
// wildcards, literals and tuples of those — no constructors, because
// there are no datatypes to construct.
type Pat interface {
	String() string
}

// WildcardPat is the _ pattern.
type WildcardPat struct{}

// VarPat is a variable pattern.
type VarPat struct {
	Name string
}

// ConstPat is an integer or boolean literal pattern.
type ConstPat struct {
	Text string
}

// TuplePat is a tuple pattern.
type TuplePat struct {
	Items []Pat
}

func (e IntLit) String() string {
	return e.Text
}

func (e BoolLit) String() string {
	return strconv.FormatBool(e.Value)
}

func (e Ident) String() string {
	return e.Name
}

func (e If) String() string {
	return paren("if", e.Cond.String(), e.Then.String(), e.Else.String())
}

func (e Fn) String() string {
	return paren("fn", e.Arg.String(), e.Body.String())
}

func (e App) String() string {
	return paren("app", e.Fn.String(), e.Arg.String())
}

func (e Infix) String() string {
	return paren(e.Op, e.L.String(), e.R.String())
}

func (e Tuple) String() string {
	parts := make([]string, 0, len(e.Items)+1)
	parts = append(parts, "tuple")

	for _, item := range e.Items {
		parts = append(parts, item.String())
	}

	return paren(parts...)
}

func (e Let) String() string {
	parts := []string{"let"}

	for _, d := range e.Decls {
		parts = append(parts, d.String())
	}

	parts = append(parts, "in", e.Body.String(), "end")

	return paren(parts...)
}

func (d ValDecl) String() string {
	parts := []string{"val"}

	for _, b := range d.Binds {
		parts = append(parts, paren(b.Pat.String(), b.E.String()))
	}

	return paren(parts...)
}

func (d FunDecl) String() string {
	parts := []string{"fun", d.Name}

	for _, arg := range d.Args {
		parts = append(parts, arg.String())
	}

	parts = append(parts, d.Body.String())

	return paren(parts...)
}

func (WildcardPat) String() string {
	return "_"
}

func (p VarPat) String() string {
	return p.Name
}

func (p ConstPat) String() string {
	return p.Text
}

func (p TuplePat) String() string {
	parts := make([]string, 0, len(p.Items)+1)
	parts = append(parts, "tuple")

	for _, item := range p.Items {
		parts = append(parts, item.String())
	}

	return paren(parts...)
}

// paren renders parts as a parenthesized space-separated form.
func paren(parts ...string) string {
	out := ""

	for i, part := range parts {
		if i > 0 {
			out += " "
		}

		out += part
	}

	return "(" + out + ")"
}
