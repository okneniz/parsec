package tml

import (
	"errors"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/lang"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

// state is the parser state of tiny ml: the stack of the binder
// names a nested expression may reference, innermost last. The
// builtin ~ sits at the bottom of the initial stack; every fn
// pattern, every declaration and every let pushes its binders for
// the nested scopes and pops them on the way out — the state half
// of the backtracking contract, because Try rewinds only the buffer.
type state struct {
	binders []string
}

func newState() *state {
	return &state{binders: []string{"~"}}
}

func (s *state) push(names ...string) {
	s.binders = append(s.binders, names...)
}

// popTo truncates the stack back to n binders.
func (s *state) popTo(n int) {
	s.binders = s.binders[:n]
}

// lookup returns the de Bruijn index of the innermost binder named
// name: 0 is the innermost binder of all.
func (s *state) lookup(name string) (int, bool) {
	for i := len(s.binders) - 1; i >= 0; i-- {
		if s.binders[i] == name {
			return len(s.binders) - 1 - i, true
		}
	}

	return 0, false
}

// builtinVar is the reference to a builtin — the binder at the
// bottom of the stack, where ~ sits alone.
func (s *state) builtinVar() Var {
	return Var{Index: len(s.binders) - 1}
}

// parser holds the expression ladder of the language. The ladder
// itself is a lang.Parser: tiny ml declares its operators and its
// forms on it — if and fn are forms of the expression level, let,
// the negation ~ and the parenthesis are forms of the operand
// level. Every method of this type is a constructor of one part
// around the ladder, the declarations and the patterns: the
// combinators it needs are built once, above the returned closure,
// and the references that would loop the build — back into the
// ladder from a form, into atpat from tuplePat — come with the
// ladder or stay in the body.
type parser struct {
	expr parsec.Combinator[rune, Position, Expr, *state]
}

func newParser() *parser {
	p := &parser{}

	ops := lang.NewParser[*state]().
		Infix(7, "*", "/").
		Infix(6, "+", "-").
		Infix(4, "=", "<>", ">", ">=", "<", "<=").
		InfixRight(-1, "andalso").
		InfixRight(-2, "orelse").
		ExprForm(p.ifForm, p.fnForm).
		PrefixForm(p.letForm, p.negForm, p.parenForm)

	p.expr = ops.Expr(lexer,
		boolOperand("true", true),
		boolOperand("false", false),
		varOperand,
	)

	return p
}

// Parse lexes src and parses it as a tiny ml program: a sequence of
// val and fun declarations, resolved to de Bruijn indices.
func Parse(src string) ([]Decl, error) {
	buf := strings.Buffer([]rune(src))
	st := newState()

	decls, err := newParser().decList(lexer)(st, buf)
	if err != nil {
		return nil, err
	}

	fixup(decls)

	// the end check asks the stepper itself what is left: trailing
	// whitespace and comments are not input, anything else is
	// unexpected. The probe is rolled back, so the unexpected error
	// names the pending token itself.
	pos := buf.Position()

	_, lerr := lexer(st, buf)

	if seekErr := buf.Seek(pos); seekErr != nil {
		return nil, parsec.NewParseError(pos, seekErr.Error())
	}

	if lerr != nil && !errors.Is(lerr, parsec.ErrEndOfFile) {
		return nil, lerr
	}

	if lerr == nil {
		_, uerr := parsec.Unexpected(lexer, "expected declaration")(st, buf)

		return nil, uerr
	}

	return decls, nil
}

// layout lists the binder names of the top level in declaration
// order, the builtin ~ first: the row the fixup pass, the
// typechecker and the interpreter all address positionally.
func layout(decls []Decl) []string {
	names := []string{"~"}

	for _, d := range decls {
		switch decl := d.(type) {
		case ValDecl:
			for _, b := range decl.Binds {
				names = append(names, patVars(b.Pat)...)
			}
		case FunDecl:
			names = append(names, decl.Name)
		}
	}

	return names
}

// fixup resolves the forward references that name a top-level fun
// declared later — the mutual recursion of tiny ml — into TopVar
// slots of the top-level frame. A forward reference to anything else
// names no binder at all; the typechecker reports it.
func fixup(decls []Decl) {
	funs := map[string]int{}

	pos := 1
	for _, d := range decls {
		switch decl := d.(type) {
		case ValDecl:
			for _, b := range decl.Binds {
				pos += len(patVars(b.Pat))
			}
		case FunDecl:
			if _, seen := funs[decl.Name]; !seen {
				funs[decl.Name] = pos
			}

			pos++
		}
	}

	var walk func(e Expr) Expr

	walk = func(e Expr) Expr {
		switch expr := e.(type) {
		case ForwardRef:
			if pos, ok := funs[expr.Name]; ok {
				return TopVar{Pos: pos}
			}

			return expr
		case If:
			expr.Cond = walk(expr.Cond)
			expr.Then = walk(expr.Then)
			expr.Else = walk(expr.Else)

			return expr
		case Fn:
			expr.Body = walk(expr.Body)

			return expr
		case App:
			expr.Fn = walk(expr.Fn)
			expr.Arg = walk(expr.Arg)

			return expr
		case Infix:
			expr.L = walk(expr.L)
			expr.R = walk(expr.R)

			return expr
		case Tuple:
			for i, item := range expr.Items {
				expr.Items[i] = walk(item)
			}

			return expr
		case Let:
			for i := range expr.Decls {
				switch decl := expr.Decls[i].(type) {
				case ValDecl:
					for j := range decl.Binds {
						decl.Binds[j].E = walk(decl.Binds[j].E)
					}

					expr.Decls[i] = decl
				case FunDecl:
					decl.Body = walk(decl.Body)

					expr.Decls[i] = decl
				}
			}

			expr.Body = walk(expr.Body)

			return expr
		}

		return e
	}

	for i := range decls {
		switch decl := decls[i].(type) {
		case ValDecl:
			for j := range decl.Binds {
				decl.Binds[j].E = walk(decl.Binds[j].E)
			}

			decls[i] = decl
		case FunDecl:
			decl.Body = walk(decl.Body)

			decls[i] = decl
		}
	}
}

// ---------------------------------------------------------------------------
// Declarations

// decList parses declarations until the next token cannot start one.
// The caller checks its own end condition: the end of input, or the
// in keyword of a let. The leading token selects the form from a
// MapAs table keyed by the whole token and is discarded; a new
// declaration form is a new table key. A declaration that starts but
// does not parse stops the list as well — the caller's end check
// reports it.
func (p *parser) decList(
	lex tokens.Lexer[Kind, Lexeme, *state],
) parsec.Combinator[rune, Position, []Decl, *state] {
	forms := map[Token]parsec.Combinator[rune, Position, Decl, *state]{
		{Kind: KindKeyword, Lexeme: "val"}: p.valDecl(lex),
		{Kind: KindKeyword, Lexeme: "fun"}: p.funDecl(lex),
	}

	decl := parsec.MapAs("expected declaration", forms, tokens.Satisfy(lex, "expected declaration", parsec.Anything[Token]))

	return parsec.Many(0, parsec.Try(decl))
}

// valDecl parses the rest of a val declaration, pat = exp; the val
// keyword is already consumed. The expression is read in the scope
// the declaration appears in; its pattern binds over everything that
// follows, so the binders are pushed only after it.
func (p *parser) valDecl(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Decl, *state] {
	eq := tokens.Exact(lex, KindOperator, "=")
	atpat := p.atpat(lex)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Decl, parsec.Error[Position]) {
		pat, perr := atpat(st, buf)
		if perr != nil {
			return nil, perr
		}

		if _, err := eq(st, buf); err != nil {
			return nil, err
		}

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		e, eerr := p.expr(st, buf)
		if eerr != nil {
			return nil, eerr
		}

		st.push(patVars(pat)...)

		return ValDecl{Binds: []ValBind{{Pat: pat, E: e}}}, nil
	}
}

// funDecl parses the rest of a recursive function declaration, fun
// f x1 ... xn = exp with atomic argument patterns; the fun keyword
// is already consumed. The name and the arguments bind inside the
// body — recursion — and, like every declaration binder, over
// everything that follows.
func (p *parser) funDecl(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Decl, *state] {
	eq := tokens.Exact(lex, KindOperator, "=")
	ident := tokens.OfKind(lex, KindIdent)
	args := parsec.Many(0, parsec.Try(p.atpat(lex)))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Decl, parsec.Error[Position]) {
		name, nerr := ident(st, buf)
		if nerr != nil {
			return nil, nerr
		}

		decl := FunDecl{Name: string(name.Lexeme)}

		// Many never fails: it stops at the first argument pattern
		// that does not parse
		decl.Args, _ = args(st, buf)

		if _, err := eq(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()

		st.push(decl.Name)
		st.push(patVarsAll(decl.Args)...)

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		body, berr := p.expr(st, buf)
		if berr != nil {
			st.popTo(base)

			return nil, berr
		}

		// the arguments scope over the body only; the name binds
		// over the rest of the program
		st.popTo(base + 1)

		decl.Body = body

		return decl, nil
	}
}

// depth returns the number of binders on the stack.
func (s *state) depth() int {
	return len(s.binders)
}

// ---------------------------------------------------------------------------
// Expression forms

// ifForm builds the if-expression form, if exp then exp else exp. It
// is a form of the expression level — an if is not an operand, and
// every branch is again a full expression.
func (p *parser) ifForm(l lang.Ladder[*state]) parsec.Combinator[rune, Position, Expr, *state] {
	ifKw := tokens.Exact(l.Lex, KindKeyword, "if")
	thenKw := tokens.Exact(l.Lex, KindKeyword, "then")
	elseKw := tokens.Exact(l.Lex, KindKeyword, "else")

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := ifKw(st, buf); err != nil {
			return nil, err
		}

		cond, cerr := l.Expr(st, buf)
		if cerr != nil {
			return nil, cerr
		}

		if _, err := thenKw(st, buf); err != nil {
			return nil, err
		}

		then, terr := l.Expr(st, buf)
		if terr != nil {
			return nil, terr
		}

		if _, err := elseKw(st, buf); err != nil {
			return nil, err
		}

		els, eerr := l.Expr(st, buf)
		if eerr != nil {
			return nil, eerr
		}

		return If{Cond: cond, Then: then, Else: els}, nil
	}
}

// fnForm builds the function form, fn pat => exp. It is a form of
// the expression level, like if: a lambda is not an operand. The
// pattern binds over the body — pushed before it and popped after,
// on the failure paths too, because a retried alternative must see
// the stack as it was.
func (p *parser) fnForm(l lang.Ladder[*state]) parsec.Combinator[rune, Position, Expr, *state] {
	fnKw := tokens.Exact(l.Lex, KindKeyword, "fn")
	arrow := tokens.Exact(l.Lex, KindOperator, "=>")
	atpat := p.atpat(l.Lex)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := fnKw(st, buf); err != nil {
			return nil, err
		}

		arg, aerr := atpat(st, buf)
		if aerr != nil {
			return nil, aerr
		}

		if _, err := arrow(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()

		st.push(patVars(arg)...)
		defer st.popTo(base)

		body, berr := l.Expr(st, buf)
		if berr != nil {
			return nil, berr
		}

		return Fn{Arg: arg, Body: body}, nil
	}
}

// letForm builds the let-expression form, decl+ in exp end. It is
// a form of the operand level — a let is an atom, f let ... end applies
// f — and its body is a full expression. The declarations bind
// sequentially over the body; the whole block pops its binders on
// the way out, on the failure paths too.
func (p *parser) letForm(l lang.Ladder[*state]) parsec.Combinator[rune, Position, Expr, *state] {
	letKw := tokens.Exact(l.Lex, KindKeyword, "let")
	inKw := tokens.Exact(l.Lex, KindKeyword, "in")
	endKw := tokens.Exact(l.Lex, KindKeyword, "end")
	decs := p.decList(l.Lex)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := letKw(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()
		defer st.popTo(base)

		decls, derr := decs(st, buf)
		if derr != nil {
			return nil, derr
		}

		if len(decls) == 0 {
			return nil, parsec.NewParseError(buf.Position(), "expected declaration in let")
		}

		if _, err := inKw(st, buf); err != nil {
			return nil, err
		}

		body, berr := l.Expr(st, buf)
		if berr != nil {
			return nil, berr
		}

		if _, err := endKw(st, buf); err != nil {
			return nil, err
		}

		return Let{Decls: decls, Body: body}, nil
	}
}

// negForm builds the negated literal ~3, with the negation glued to
// the literal so that f ~3 is f applied to ~3, not f ~ applied to 3.
// A lone ~ is the negation function itself — the builtin binder of
// the initial stack.
func (p *parser) negForm(l lang.Ladder[*state]) parsec.Combinator[rune, Position, Expr, *state] {
	tilde := tokens.Exact(l.Lex, KindOperator, "~")
	intLit := parsec.Try(tokens.OfKind(l.Lex, KindInt))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := tilde(st, buf); err != nil {
			return nil, err
		}

		neg := st.builtinVar()

		lit, lerr := intLit(st, buf)
		if lerr != nil {
			return neg, nil
		}

		value, _ := strconv.ParseInt(string(lit.Lexeme), 10, 64)

		return App{Fn: neg, Arg: IntLit{Text: string(lit.Lexeme), Value: value}}, nil
	}
}

// parenForm builds the grouping form (e) and the tuple form
// (e1, ..., en); every element is again a full expression.
func (p *parser) parenForm(l lang.Ladder[*state]) parsec.Combinator[rune, Position, Expr, *state] {
	open := tokens.Exact(l.Lex, KindSymbol, "(")
	closeParen := tokens.Exact(l.Lex, KindSymbol, ")")
	comma := parsec.Try(tokens.Exact(l.Lex, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := open(st, buf); err != nil {
			return nil, err
		}

		first, ferr := l.Expr(st, buf)
		if ferr != nil {
			return nil, ferr
		}

		items := []Expr{first}

		for {
			if _, cerr := comma(st, buf); cerr != nil {
				break
			}

			e, eerr := l.Expr(st, buf)
			if eerr != nil {
				return nil, eerr
			}

			items = append(items, e)
		}

		if _, err := closeParen(st, buf); err != nil {
			return nil, err
		}

		if len(items) == 1 {
			return items[0], nil
		}

		return Tuple{Items: items}, nil
	}
}

// boolOperand turns a boolean literal into an extra operand of the
// ladder: a keyword with a fixed meaning needs no recursion.
func boolOperand(lexeme Lexeme, value bool) lang.Operand[*state] {
	return func(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Expr, *state] {
		return parsec.Cast(
			tokens.Exact(lex, KindKeyword, lexeme),
			func(Token) (Expr, error) {
				return BoolLit{Value: value}, nil
			},
		)
	}
}

// varOperand is the extra operand that resolves identifiers into
// variable references: a name on the binder stack becomes its de
// Bruijn index, a name off the stack a forward reference the fixup
// pass of Parse resolves against the whole program.
func varOperand(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Expr, *state] {
	ident := tokens.OfKind(lex, KindIdent)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		t, err := ident(st, buf)
		if err != nil {
			return nil, err
		}

		name := string(t.Lexeme)

		if i, bound := st.lookup(name); bound {
			return Var{Index: i}, nil
		}

		return ForwardRef{Name: name}, nil
	}
}

// ---------------------------------------------------------------------------
// Patterns

// atpat parses one atomic pattern. Tiny ml has no constructors, so
// there are no infix constructor patterns and no layered patterns.
// The dispatch mirrors the operand level of the ladder: a Choice of
// Try-wrapped alternatives.
func (p *parser) atpat(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Pat, *state] {
	wildcard := parsec.Cast(tokens.Exact(lex, KindSymbol, "_"), func(Token) (Pat, error) {
		return WildcardPat{}, nil
	})

	constPat := parsec.Cast(
		tokens.Satisfy(lex, "expected literal pattern", isLiteralToken),
		func(t Token) (Pat, error) {
			return ConstPat{Text: string(t.Lexeme)}, nil
		},
	)

	varPat := parsec.Cast(
		tokens.OfKind(lex, KindIdent),
		func(t Token) (Pat, error) {
			return VarPat{Name: string(t.Lexeme)}, nil
		},
	)

	return parsec.Choice("expected pattern",
		parsec.Try(wildcard),
		parsec.Try(constPat),
		parsec.Try(p.negPat(lex)),
		parsec.Try(varPat),
		parsec.Try(p.tuplePat(lex)),
	)
}

// isBoolToken reports whether the token is a boolean literal.
func isBoolToken(t Token) bool {
	return t.Kind == KindKeyword && (t.Lexeme == "true" || t.Lexeme == "false")
}

// isLiteralToken reports whether the token is an integer or boolean
// literal.
func isLiteralToken(t Token) bool {
	return t.Kind == KindInt || isBoolToken(t)
}

// negPat parses the negative literal pattern ~1; a lone ~ is a
// variable pattern.
func (p *parser) negPat(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Pat, *state] {
	tilde := tokens.Exact(lex, KindOperator, "~")
	intLit := parsec.Try(tokens.OfKind(lex, KindInt))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Pat, parsec.Error[Position]) {
		if _, err := tilde(st, buf); err != nil {
			return nil, err
		}

		lit, lerr := intLit(st, buf)
		if lerr != nil {
			return VarPat{Name: "~"}, nil
		}

		return ConstPat{Text: "~" + string(lit.Lexeme)}, nil
	}
}

// tuplePat parses (p) and tuple patterns (p1, ..., pn).
func (p *parser) tuplePat(lex tokens.Lexer[Kind, Lexeme, *state]) parsec.Combinator[rune, Position, Pat, *state] {
	open := tokens.Exact(lex, KindSymbol, "(")
	closeParen := tokens.Exact(lex, KindSymbol, ")")
	comma := parsec.Try(tokens.Exact(lex, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Pat, parsec.Error[Position]) {
		if _, err := open(st, buf); err != nil {
			return nil, err
		}

		// atpat's Choice builds this form, so the reference stays
		// in the body; one build serves every element
		atpat := p.atpat(lex)

		var items []Pat

		for {
			pat, perr := atpat(st, buf)
			if perr != nil {
				return nil, perr
			}

			items = append(items, pat)

			if _, cerr := comma(st, buf); cerr != nil {
				break
			}
		}

		if _, err := closeParen(st, buf); err != nil {
			return nil, err
		}

		if len(items) == 1 {
			return items[0], nil
		}

		return TuplePat{Items: items}, nil
	}
}

// patVars lists the binder names of a pattern in binding order:
// pushed left to right, the rightmost ends up innermost.
func patVars(p Pat) []string {
	switch pat := p.(type) {
	case VarPat:
		return []string{pat.Name}
	case TuplePat:
		return patVarsAll(pat.Items)
	}

	return nil
}

func patVarsAll(pats []Pat) []string {
	var out []string

	for _, p := range pats {
		out = append(out, patVars(p)...)
	}

	return out
}
