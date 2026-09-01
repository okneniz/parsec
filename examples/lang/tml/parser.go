package tml

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

// fixity is the infix status of an operator: a precedence level and
// an associativity. The andalso and orelse keywords live in the same
// space, below every operator level.
type fixity struct {
	prec  int
	right bool
}

// defaultFixity is the one and only fixity table of tiny ml: the
// language has no infix declarations, so the table is static.
func defaultFixity() map[Lexeme]fixity {
	return map[Lexeme]fixity{
		"*": {7, false}, "/": {7, false},
		"+": {6, false}, "-": {6, false},
		"=": {4, false}, "<>": {4, false}, ">": {4, false},
		">=": {4, false}, "<": {4, false}, "<=": {4, false},
		"andalso": {-1, true},
		"orelse":  {-2, true},
	}
}

// parser holds the fixity table. Every method is a constructor: the
// combinators it needs are built once, above the returned closure,
// and the closure captures them. The exception is a reference that
// would loop the build — back into expr from ifExp, fnExp, parenExp,
// and letExp, into atpat from tuplePat, into decList from letExp —
// those stay in the body, built once per invocation.
type parser struct {
	fixity map[Lexeme]fixity
}

func newParser() *parser {
	return &parser{fixity: defaultFixity()}
}

// Parse lexes src and parses it as a tiny ml program: a sequence of
// val and fun declarations.
func Parse(src string) ([]Decl, error) {
	buf := strings.Buffer([]rune(src))

	decls, err := newParser().decList(lexer)(buf)
	if err != nil {
		return nil, err
	}

	// the end check asks the stepper itself what is left: trailing
	// whitespace and comments are not input, anything else is
	// unexpected. The probe is rolled back, so the unexpected error
	// names the pending token itself.
	pos := buf.Position()

	_, lerr := lexer(buf)

	if seekErr := buf.Seek(pos); seekErr != nil {
		return nil, parsec.NewParseError(pos, seekErr.Error())
	}

	if lerr != nil && !errors.Is(lerr, parsec.ErrEndOfFile) {
		return nil, lerr
	}

	if lerr == nil {
		_, uerr := parsec.Unexpected(lexer, "expected declaration")(buf)

		return nil, uerr
	}

	return decls, nil
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
	lex tokens.Lexer[Kind, Lexeme],
) parsec.Combinator[rune, Position, []Decl] {
	forms := map[Token]parsec.Combinator[rune, Position, Decl]{
		{Kind: KindKeyword, Lexeme: "val"}: p.valDecl(lex),
		{Kind: KindKeyword, Lexeme: "fun"}: p.funDecl(lex),
	}

	decl := parsec.MapAs("expected declaration", forms, tokens.Satisfy(lex, "expected declaration", parsec.Anything[Token]))

	return parsec.Many(0, parsec.Try(decl))
}

// valDecl parses the rest of a val declaration, pat = exp; the val
// keyword is already consumed.
func (p *parser) valDecl(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Decl] {
	eq := tokens.Exact(lex, KindOperator, "=")
	atpat := p.atpat(lex)
	expr := p.expr(lex)

	return func(buf parsec.Buffer[rune, Position]) (Decl, parsec.Error[Position]) {
		pat, perr := atpat(buf)
		if perr != nil {
			return nil, perr
		}

		if _, err := eq(buf); err != nil {
			return nil, err
		}

		e, eerr := expr(buf)
		if eerr != nil {
			return nil, eerr
		}

		return ValDecl{Binds: []ValBind{{Pat: pat, E: e}}}, nil
	}
}

// funDecl parses the rest of a recursive function declaration, fun
// f x1 ... xn = exp with atomic argument patterns; the fun keyword
// is already consumed.
func (p *parser) funDecl(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Decl] {
	eq := tokens.Exact(lex, KindOperator, "=")
	ident := tokens.OfKind(lex, KindIdent)
	args := parsec.Many(0, parsec.Try(p.atpat(lex)))
	expr := p.expr(lex)

	return func(buf parsec.Buffer[rune, Position]) (Decl, parsec.Error[Position]) {
		name, nerr := ident(buf)
		if nerr != nil {
			return nil, nerr
		}

		decl := FunDecl{Name: string(name.Lexeme)}

		// Many never fails: it stops at the first argument pattern
		// that does not parse
		decl.Args, _ = args(buf)

		if _, err := eq(buf); err != nil {
			return nil, err
		}

		body, berr := expr(buf)
		if berr != nil {
			return nil, berr
		}
		decl.Body = body

		return decl, nil
	}
}

// ---------------------------------------------------------------------------
// Expressions

// expr parses a full expression: one of the special forms if and fn,
// or an infix chain. The special forms dispatch through a MapAs table
// keyed by the whole leading token — a discarded prefix, the form
// owns the input from there; anything else is the infix chain. A new
// special form is a new table key.
func (p *parser) expr(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	token := tokens.Satisfy(
		lex,
		"expected expression",
		parsec.Anything[Token],
	)

	form := parsec.MapAs(
		"expected expression",
		map[Token]parsec.Combinator[rune, Position, Expr]{
			{Kind: KindKeyword, Lexeme: "if"}: p.ifExp(lex),
			{Kind: KindKeyword, Lexeme: "fn"}: p.fnExp(lex),
		},
		token,
	)

	return parsec.Choice("expected expression",
		parsec.Try(form),
		p.infexp(lex),
	)
}

// ifExp parses the rest of an if-expression, exp then exp else exp;
// the if keyword is already consumed.
func (p *parser) ifExp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	thenKw := tokens.Exact(lex, KindKeyword, "then")
	elseKw := tokens.Exact(lex, KindKeyword, "else")

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		// expr's table builds this form, so the reference stays in
		// the body; one build serves all three branches
		expr := p.expr(lex)

		cond, cerr := expr(buf)
		if cerr != nil {
			return nil, cerr
		}

		if _, err := thenKw(buf); err != nil {
			return nil, err
		}

		then, terr := expr(buf)
		if terr != nil {
			return nil, terr
		}

		if _, err := elseKw(buf); err != nil {
			return nil, err
		}

		els, eerr := expr(buf)
		if eerr != nil {
			return nil, eerr
		}

		return If{Cond: cond, Then: then, Else: els}, nil
	}
}

// fnExp parses the rest of an fn-expression, pat => exp; the fn
// keyword is already consumed.
func (p *parser) fnExp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	arrow := tokens.Exact(lex, KindOperator, "=>")
	atpat := p.atpat(lex)

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		arg, aerr := atpat(buf)
		if aerr != nil {
			return nil, aerr
		}

		if _, err := arrow(buf); err != nil {
			return nil, err
		}

		// the reference back into expr's table cannot move above
		// the closure
		body, berr := p.expr(lex)(buf)
		if berr != nil {
			return nil, berr
		}

		return Fn{Arg: arg, Body: body}, nil
	}
}

// orItem is one element of a flat infix sequence: an operand, or an
// operator when op is set.
type orItem struct {
	op   Lexeme
	expr Expr
}

// resolveInfix folds a flat operand-operator sequence into a tree with
// the shunting-yard algorithm, using the fixity table for precedence
// and associativity.
func resolveInfix(items []orItem, fix map[Lexeme]fixity) Expr {
	var vals []Expr
	var ops []Lexeme

	apply := func() {
		r := vals[len(vals)-1]
		l := vals[len(vals)-2]
		vals = vals[:len(vals)-2]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		vals = append(vals, Infix{Op: string(op), L: l, R: r})
	}

	for _, it := range items {
		if it.op == "" {
			vals = append(vals, it.expr)
			continue
		}

		cur := fix[it.op]
		for len(ops) > 0 {
			top := fix[ops[len(ops)-1]]
			if top.prec > cur.prec || top.prec == cur.prec && !cur.right {
				apply()
			} else {
				break
			}
		}

		ops = append(ops, it.op)
	}

	for len(ops) > 0 {
		apply()
	}

	return vals[0]
}

// infexp parses a flat sequence of applications separated by infix
// operators and resolves it with the fixity table.
func (p *parser) infexp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	op := parsec.Try(tokens.Satisfy(lex, "operator", p.isOpToken))
	appexp := p.appexp(lex)

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		first, err := appexp(buf)
		if err != nil {
			return nil, err
		}

		items := []orItem{{expr: first}}

		for {
			next, oerr := op(buf)
			if oerr != nil {
				break
			}

			rhs, rerr := appexp(buf)
			if rerr != nil {
				return nil, parsec.NewParseError(
					buf.Position(),
					fmt.Sprintf("expected operand after operator %s", next.Lexeme),
				)
			}

			items = append(items, orItem{op: next.Lexeme}, orItem{expr: rhs})
		}

		return resolveInfix(items, p.fixity), nil
	}
}

// isOpToken reports whether the token acts as an infix operator: an
// operator token with fixity status, or the andalso and orelse
// keywords, which live in the same space below every operator level.
func (p *parser) isOpToken(t Token) bool {
	if t.Kind == KindKeyword {
		return t.Lexeme == "andalso" || t.Lexeme == "orelse"
	}

	_, hasFixity := p.fixity[t.Lexeme]
	return t.Kind == KindOperator && hasFixity
}

// appexp parses left-nested application: one atom, then as many
// argument atoms as parse, each attempt backtracked when it fails.
// Application binds tighter than any infix operator, and the infix
// operators are their own token kind — f x + y is (f x) + y by the
// dispatch, not by a predicate.
func (p *parser) appexp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	first := p.atexp(lex)
	args := parsec.Many(0, parsec.Try(first))

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		e, err := first(buf)
		if err != nil {
			return nil, err
		}

		rest, _ := args(buf)

		for _, arg := range rest {
			e = App{Fn: e, Arg: arg}
		}

		return e, nil
	}
}

// atexp parses one atomic expression. The dispatch is a cascade:
// the forms whose whole leading token identifies them — let, the
// boolean literals, the negation ~, the opening parenthesis — are
// looked up in a MapAs table keyed by the token itself, a Token is
// comparable; the token is a discarded prefix, the table entry owns
// the rest. What is left has no such prefix — an integer has no
// fixed lexeme — and stays a Choice of Try-wrapped alternatives, the
// idiom of lang.primary. A new form with a fixed leading token is a
// new table key.
func (p *parser) atexp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	exactForms := map[Token]parsec.Combinator[rune, Position, Expr]{
		{Kind: KindKeyword, Lexeme: "let"}:   p.letExp(lex),
		{Kind: KindKeyword, Lexeme: "true"}:  parsec.Const[rune, Position, Expr](BoolLit{Value: true}),
		{Kind: KindKeyword, Lexeme: "false"}: parsec.Const[rune, Position, Expr](BoolLit{Value: false}),
		{Kind: KindOperator, Lexeme: "~"}:    p.negExp(lex),
		{Kind: KindSymbol, Lexeme: "("}:      p.parenExp(lex),
	}

	token := tokens.Satisfy(lex, "expected expression", parsec.Anything[Token])
	form := parsec.MapAs("expected expression", exactForms, token)

	intLit := parsec.Cast(tokens.OfKind(lex, KindInt), func(t Token) (Expr, error) {
		value, _ := strconv.ParseInt(string(t.Lexeme), 10, 64)
		return IntLit{Text: string(t.Lexeme), Value: value}, nil
	})

	ident := parsec.Cast(
		tokens.OfKind(lex, KindIdent),
		func(t Token) (Expr, error) {
			return Ident{Name: string(t.Lexeme)}, nil
		},
	)

	return parsec.Choice("expected expression",
		parsec.Try(form),
		parsec.Try(intLit),
		parsec.Try(ident),
	)
}

// isBoolToken reports whether the token is a boolean literal.
func isBoolToken(t Token) bool {
	return t.Kind == KindKeyword && (t.Lexeme == "true" || t.Lexeme == "false")
}

// negExp parses the rest of a negated literal, ~3, with the negation
// glued to the literal so that f ~3 is f applied to ~3, not f ~
// applied to 3. A lone ~ is the negation function itself. The tilde
// is already consumed by the dispatch table.
func (p *parser) negExp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	intLit := parsec.Try(tokens.OfKind(lex, KindInt))

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		lit, lerr := intLit(buf)
		if lerr != nil {
			return Ident{Name: "~"}, nil
		}

		value, _ := strconv.ParseInt(string(lit.Lexeme), 10, 64)

		return App{Fn: Ident{Name: "~"}, Arg: IntLit{Text: string(lit.Lexeme), Value: value}}, nil
	}
}

// parenExp parses the rest of a parenthesized expression, (e) and
// tuple expressions (e1, ..., en); the opening parenthesis is
// already consumed by the dispatch table.
func (p *parser) parenExp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	closeParen := tokens.Exact(lex, KindSymbol, ")")
	comma := parsec.Try(tokens.Exact(lex, KindSymbol, ","))

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		// atexp's table builds this form under expr's subtree, so
		// the reference stays in the body; one build serves the
		// whole tuple
		expr := p.expr(lex)

		first, ferr := expr(buf)
		if ferr != nil {
			return nil, ferr
		}

		items := []Expr{first}

		for {
			if _, cerr := comma(buf); cerr != nil {
				break
			}

			e, eerr := expr(buf)
			if eerr != nil {
				return nil, eerr
			}

			items = append(items, e)
		}

		if _, err := closeParen(buf); err != nil {
			return nil, err
		}

		if len(items) == 1 {
			return items[0], nil
		}

		return Tuple{Items: items}, nil
	}
}

// letExp parses the rest of a let-expression, decl+ in exp end; the
// let keyword is already consumed by the dispatch table.
func (p *parser) letExp(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
	inKw := tokens.Exact(lex, KindKeyword, "in")
	endKw := tokens.Exact(lex, KindKeyword, "end")

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		// both references sit under atexp's table inside expr's
		// subtree and cannot move above the closure
		decls, derr := p.decList(lex)(buf)
		if derr != nil {
			return nil, derr
		}

		if len(decls) == 0 {
			return nil, parsec.NewParseError(buf.Position(), "expected declaration in let")
		}

		if _, err := inKw(buf); err != nil {
			return nil, err
		}

		body, berr := p.expr(lex)(buf)
		if berr != nil {
			return nil, berr
		}

		if _, err := endKw(buf); err != nil {
			return nil, err
		}

		return Let{Decls: decls, Body: body}, nil
	}
}

// ---------------------------------------------------------------------------
// Patterns

// atpat parses one atomic pattern. Tiny ml has no constructors, so
// there are no infix constructor patterns and no layered patterns.
// The dispatch mirrors atexp: a Choice of Try-wrapped alternatives.
func (p *parser) atpat(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Pat] {
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

// isLiteralToken reports whether the token is an integer or boolean
// literal.
func isLiteralToken(t Token) bool {
	return t.Kind == KindInt || isBoolToken(t)
}

// negPat parses the negative literal pattern ~1; a lone ~ is a
// variable pattern.
func (p *parser) negPat(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Pat] {
	tilde := tokens.Exact(lex, KindOperator, "~")
	intLit := parsec.Try(tokens.OfKind(lex, KindInt))

	return func(buf parsec.Buffer[rune, Position]) (Pat, parsec.Error[Position]) {
		if _, err := tilde(buf); err != nil {
			return nil, err
		}

		lit, lerr := intLit(buf)
		if lerr != nil {
			return VarPat{Name: "~"}, nil
		}

		return ConstPat{Text: "~" + string(lit.Lexeme)}, nil
	}
}

// tuplePat parses (p) and tuple patterns (p1, ..., pn).
func (p *parser) tuplePat(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Pat] {
	open := tokens.Exact(lex, KindSymbol, "(")
	closeParen := tokens.Exact(lex, KindSymbol, ")")
	comma := parsec.Try(tokens.Exact(lex, KindSymbol, ","))

	return func(buf parsec.Buffer[rune, Position]) (Pat, parsec.Error[Position]) {
		if _, err := open(buf); err != nil {
			return nil, err
		}

		// atpat's Choice builds this form, so the reference stays
		// in the body; one build serves every element
		atpat := p.atpat(lex)

		var items []Pat

		for {
			pat, perr := atpat(buf)
			if perr != nil {
				return nil, perr
			}

			items = append(items, pat)

			if _, cerr := comma(buf); cerr != nil {
				break
			}
		}

		if _, err := closeParen(buf); err != nil {
			return nil, err
		}

		if len(items) == 1 {
			return items[0], nil
		}

		return TuplePat{Items: items}, nil
	}
}
