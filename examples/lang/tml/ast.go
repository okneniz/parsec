package tml

import (
	"strconv"

	"github.com/okneniz/parsec/lang"
)

// Decl is a declaration node.
type Decl interface {
	String() string
}

// Expr is an expression node. The contract is the one of the lang
// ladder: a node renders itself as a string.
type Expr = lang.Expr

// The nodes the expression ladder of lang produces are its own: the
// rendered forms of tiny ml match, so the ladder and the typechecker
// and the interpreter share them as is.
type IntLit = lang.IntLit

// BoolLit is a boolean literal.
type BoolLit struct {
	Value bool
}

func (e BoolLit) String() string {
	return strconv.FormatBool(e.Value)
}

// Var is a variable reference by de Bruijn index: 0 is the innermost
// enclosing binder, 1 the next one out.
type Var struct {
	Index int
}

func (e Var) String() string {
	return "var(" + strconv.Itoa(e.Index) + ")"
}

// FreeVar is a free variable: a reference to a name no enclosing
// binder holds when it is read. Nothing resolves it later — the
// runtime knows indices only, so the typechecker reports it.
type FreeVar struct {
	Name string
}

func (e FreeVar) String() string {
	return e.Name
}

// If is a conditional expression.
type If struct {
	Cond, Then, Else Expr
}

func (e If) String() string {
	return paren("if", e.Cond.String(), e.Then.String(), e.Else.String())
}

// Fn is a single-clause function expression: fn pat => body.
type Fn struct {
	Arg  Pat
	Body Expr
}

func (e Fn) String() string {
	return paren("fn", e.Body.String())
}

// App is function application; application binds tighter than any
// infix operator.
type App = lang.App

// Infix is an infix application; Op is the operator lexeme, which may
// also be the reserved words andalso and orelse.
type Infix = lang.Binary

// Tuple is a tuple expression (e1, ..., en).
type Tuple struct {
	Items []Expr
}

func (e Tuple) String() string {
	if len(e.Items) == 0 {
		return "()"
	}

	parts := make([]string, 0, len(e.Items)+1)
	parts = append(parts, "tuple")

	for _, item := range e.Items {
		parts = append(parts, item.String())
	}

	return paren(parts...)
}

// Let is a let expression with a single-expression body.
type Let struct {
	Decls []Decl
	Body  Expr
}

func (e Let) String() string {
	parts := []string{"let"}

	for _, d := range e.Decls {
		parts = append(parts, d.String())
	}

	parts = append(parts, "in", e.Body.String(), "end")

	return paren(parts...)
}

// ValDecl is a val declaration.
type ValDecl struct {
	Binds []ValBind
}

// ValBind is one pat = exp binding of a val declaration.
type ValBind struct {
	Pat Pat
	E   Expr
}

func (d ValDecl) String() string {
	parts := []string{"val"}

	for _, b := range d.Binds {
		parts = append(parts, paren(b.Pat.String(), b.E.String()))
	}

	return paren(parts...)
}

// FunDecl is a recursive function declaration: fun f x1 ... xn = body,
// with atomic argument patterns.
type FunDecl struct {
	Name string
	Args []Pat
	Body Expr
}

func (d FunDecl) String() string {
	parts := []string{"fun", d.Name}

	for _, arg := range d.Args {
		parts = append(parts, arg.String())
	}

	parts = append(parts, d.Body.String())

	return paren(parts...)
}

// Pat is a pattern node. Tiny ml patterns are first-order: variables,
// wildcards, literals and tuples of those — no constructors, because
// there are no datatypes to construct.
type Pat interface {
	String() string
}

// WildcardPat is the _ pattern.
type WildcardPat struct{}

func (WildcardPat) String() string {
	return "_"
}

// VarPat is a variable pattern.
type VarPat struct {
	Name string // TODO : здесь должен быть индекс
}

func (p VarPat) String() string {
	return p.Name
}

// ConstPat is an integer or boolean literal pattern.
type ConstPat struct {
	Text string
}

func (p ConstPat) String() string {
	return p.Text
}

// TuplePat is a tuple pattern.
type TuplePat struct {
	Items []Pat
}

func (p TuplePat) String() string {
	if len(p.Items) == 0 {
		return "()"
	}

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
