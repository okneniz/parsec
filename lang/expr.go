package lang

import (
	"math"
	"slices"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

type Expr interface {
	String() string
}

type IntLit struct {
	Text  string
	Value int64
}

type Ident struct {
	Name string
}

type App struct {
	Fn  Expr
	Arg Expr
}

func (e App) String() string {
	return "(app " + e.Fn.String() + " " + e.Arg.String() + ")"
}

type Unary struct {
	Op  string
	Sub Expr
}

type Postfix struct {
	Op string
	E  Expr
}

type Binary struct {
	Op   string
	L, R Expr
}

type Conditional struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (e IntLit) String() string {
	return e.Text
}

func (e Ident) String() string {
	return e.Name
}

func (e Unary) String() string {
	return "(" + e.Op + " " + e.Sub.String() + ")"
}

func (e Postfix) String() string {
	return "(" + e.Op + " " + e.E.String() + ")"
}

func (e Binary) String() string {
	return "(" + e.Op + " " + e.L.String() + " " + e.R.String() + ")"
}

func (e Conditional) String() string {
	return "(?: " + e.Cond.String() + " " + e.Then.String() + " " + e.Else.String() + ")"
}

type Fixity struct {
	Prec  int
	Right bool
}

// Parser holds the operator table of a language and its syntactic
// forms. The type parameter S is the parser state its combinators
// carry — the same state the Definition the language lexes with
// is generic in.
type Parser[S any] struct {
	fixity      map[string]Fixity
	ternary     map[string]ternaryDef
	prefix      map[string]struct{}
	postfix     map[string]struct{}
	exprForms   []Form[S]
	prefixForms []Form[S]
}

// NewParser creates a parser with an empty operator table.
func NewParser[S any]() *Parser[S] {
	return &Parser[S]{
		fixity:  map[string]Fixity{},
		ternary: map[string]ternaryDef{},
		prefix:  map[string]struct{}{},
		postfix: map[string]struct{}{},
	}
}

// Infix declares left-associative operators of the given precedence;
// higher binds tighter.
func (p *Parser[S]) Infix(prec int, ops ...string) *Parser[S] {
	for _, op := range ops {
		p.fixity[op] = Fixity{Prec: prec}
	}
	return p
}

// InfixRight declares right-associative operators of the given
// precedence, as exponentiation usually is.
func (p *Parser[S]) InfixRight(prec int, ops ...string) *Parser[S] {
	for _, op := range ops {
		p.fixity[op] = Fixity{Prec: prec, Right: true}
	}
	return p
}

// Ternary declares a right-associative conditional operator of the
// given precedence, the question and the colon its two lexemes: a
// ? b : c ? d : e is a ? b : (c ? d : e), the ternary conditional
// operator of c between the assignment and the logical or. The
// condition reduces the operators tighter than the conditional into
// itself, the middle part is a full expression — the comma and the
// assignment of c live in it — and the else branch recurses at the
// precedence of the conditional. The colon must not carry fixity of
// its own: the middle part would take it for an infix operator.
func (p *Parser[S]) Ternary(prec int, question, colon string) *Parser[S] {
	p.ternary[question] = ternaryDef{colon: colon, prec: prec}
	return p
}

// Prefix declares prefix operators. A prefix operator applies to the
// next operand only, so -a + b is (+ (- a) b) and -(a + b) needs
// parentheses.
func (p *Parser[S]) Prefix(ops ...string) *Parser[S] {
	for _, op := range ops {
		p.prefix[op] = parsec.Stateless{}
	}
	return p
}

// Postfix declares postfix operators. A postfix operator applies to
// the complete operand before it and chains right, so a! + b is
// (+ (! a) b) and a!! is (! (! a)). An operator must not be both
// postfix and infix: the postfix loop would consume it first.
func (p *Parser[S]) Postfix(ops ...string) *Parser[S] {
	for _, op := range ops {
		p.postfix[op] = parsec.Stateless{}
	}
	return p
}

// ExprForm declares forms of the expression level, tried before the
// infix chain: the if e then e else e and fn p => e of sml live
// here, above every operator. An expression form is not an operand —
// 1 + if e then e else e does not parse — and it consumes its whole
// rest itself, subexpressions included, through the ladder.
func (p *Parser[S]) ExprForm(forms ...Form[S]) *Parser[S] {
	p.exprForms = append(p.exprForms, forms...)

	return p
}

// PrefixForm declares forms of the operand level, tried among the
// operands before the built-in ones: a sizeof e reads its operand
// through the ladder's Atom and so binds tighter than any infix
// operator, a cast (t) e likewise; a form that owns its
// subexpressions wholly — a let, a parenthesized tuple — parses them
// through the ladder's Expr.
func (p *Parser[S]) PrefixForm(forms ...Form[S]) *Parser[S] {
	p.prefixForms = append(p.prefixForms, forms...)

	return p
}

// Ladder hands a form the recursive entry points of the expression
// machinery, so the form parses its subexpressions with the same
// parser it is part of. Expr is the full expression, forms included;
// Atom is one operand, with its prefix and postfix operators — the
// level a prefix operator applies to; From reads the expression from
// a precedence up, the chain cut below it — the level the arguments
// of a call read at, above a comma operator. Lex is the lexer the
// whole ladder reads through. Every entry point carries the same
// parser state S.
type Ladder[S any] struct {
	Lex  tokens.Lexer[Kind, Lexeme, S]
	Expr parsec.Combinator[rune, strings.Position, Expr, S]
	Atom parsec.Combinator[rune, strings.Position, Expr, S]
	From func(prec int) parsec.Combinator[rune, strings.Position, Expr, S]
}

// A Form builds the combinator of a syntactic form that does not
// reduce to an operator lexeme: a keyword-led form such as sizeof e
// or if e then e else e, or a form led by a punctuation token, a
// cast (t) e, a parenthesized tuple. The forms of a Parser are built
// once, when its Expr runs; everything a form needs besides the
// ladder it builds above the returned closure.
type Form[S any] func(Ladder[S]) parsec.Combinator[rune, strings.Position, Expr, S]

type Operand[S any] func(lex tokens.Lexer[Kind, Lexeme, S]) parsec.Combinator[rune, strings.Position, Expr, S]

func (p *Parser[S]) Expr(
	lex tokens.Lexer[Kind, Lexeme, S],
	extras ...Operand[S],
) parsec.Combinator[rune, strings.Position, Expr, S] {
	steps := make([]parsec.Combinator[rune, strings.Position, Expr, S], len(extras), len(extras)+len(p.prefixForms))

	for i, extra := range extras {
		steps[i] = extra(lex)
	}

	// the forms are part of the ladder they recurse into, so the
	// ladder hands out closures that read full and atom at call
	// time, after the whole ladder is built
	ladder := Ladder[S]{Lex: lex}

	var full parsec.Combinator[rune, strings.Position, Expr, S]
	var atom parsec.Combinator[rune, strings.Position, Expr, S]

	ladder.Expr = func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		return full(state, buf)
	}

	ladder.Atom = func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		return atom(state, buf)
	}

	// the full level is the expression forms, if any, over the
	// operator chain; a sub-level of the ladder — the arguments of a
	// call above the comma — is the same composition with the chain
	// cut at a precedence: From reads it, the full expression is From
	// of the bottom
	formAlts := make([]parsec.Combinator[rune, strings.Position, Expr, S], len(p.exprForms))

	for i, form := range p.exprForms {
		formAlts[i] = parsec.Try(form(ladder))
	}

	// a form may read From while the steps it belongs to are still
	// being built, so the sub-level of a precedence is built on its
	// first read and kept
	sublevels := make(map[int]parsec.Combinator[rune, strings.Position, Expr, S])

	ladder.From = func(prec int) parsec.Combinator[rune, strings.Position, Expr, S] {
		return func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
			level, built := sublevels[prec]
			if !built {
				level = p.expr(lex, steps, ladder.Expr, prec)

				if len(formAlts) > 0 {
					level = parsec.Choice("expected expression", append(slices.Clone(formAlts), level)...)
				}

				sublevels[prec] = level
			}

			return level(state, buf)
		}
	}

	for _, form := range p.prefixForms {
		steps = append(steps, form(ladder))
	}

	full = p.expr(lex, steps, ladder.Expr, math.MinInt)

	if len(formAlts) > 0 {
		full = parsec.Choice("expected expression", append(slices.Clone(formAlts), full)...)
	}

	atom = p.atom(lex, steps, true, ladder.Expr)

	return full
}

func (p *Parser[S]) expr(
	lex tokens.Lexer[Kind, Lexeme, S],
	steps []parsec.Combinator[rune, strings.Position, Expr, S],
	rec parsec.Combinator[rune, strings.Position, Expr, S],
	minPrec int,
) parsec.Combinator[rune, strings.Position, Expr, S] {
	// an operator is a lexeme with fixity status, whatever kind the
	// lexer gave its token: the reserved words andalso and orelse of
	// sml are infix operators of their own precedence level. The
	// chain stops before an operator below minPrec — a sub-level of
	// the ladder reads its expression this way, the arguments of a
	// call above the comma
	op := parsec.Try(tokens.Satisfy(lex, "operator", func(t Token) bool {
		if f, hasFixity := p.fixity[string(t.Lexeme)]; hasFixity {
			return f.Prec >= minPrec
		}

		def, isConditional := p.ternary[string(t.Lexeme)]

		return isConditional && def.prec >= minPrec
	}))

	// the colons of the conditional operators match by lexeme too,
	// whatever kind their tokens carry
	colons := make(map[string]parsec.Combinator[rune, strings.Position, Token, S], len(p.ternary))

	for question, def := range p.ternary {
		colons[question] = tokens.Satisfy(lex, "conditional colon", func(t Token) bool {
			return string(t.Lexeme) == def.colon
		})
	}

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		first, err := p.appexp(lex, steps, rec)(state, buf)
		if err != nil {
			return nil, err
		}

		items := []infixItem{{expr: first}}

		for {
			next, oerr := op(state, buf)
			if oerr != nil {
				break
			}

			if def, isConditional := p.ternary[string(next.Lexeme)]; isConditional {
				// the condition is what the operators tighter than
				// the conditional reduce into; the looser ones stay
				// in the chain — the = of x = a ? b : c binds the
				// whole conditional
				items = reduceInfix(items, p.fixity, Fixity{Prec: def.prec})
				cond := items[len(items)-1].expr

				// the middle part is a full expression, like the one
				// in parentheses — the comma and the assignment of c
				// live here
				then, terr := rec(state, buf)
				if terr != nil {
					return nil, terr
				}

				if _, cerr := colons[string(next.Lexeme)](state, buf); cerr != nil {
					return nil, parsec.NewParseError(
						buf.Position(),
						"expected "+strconv.Quote(def.colon)+" to close the conditional operator "+strconv.Quote(string(next.Lexeme)),
					)
				}

				// the else branch recurses at the precedence of the
				// conditional, right associative
				other, eerr := p.expr(lex, steps, rec, def.prec)(state, buf)
				if eerr != nil {
					return nil, eerr
				}

				items[len(items)-1] = infixItem{expr: Conditional{Cond: cond, Then: then, Else: other}}

				continue
			}

			rhs, rerr := p.appexp(lex, steps, rec)(state, buf)
			if rerr != nil {
				return nil, parsec.NewParseError(
					buf.Position(),
					"expected operand after operator "+strconv.Quote(string(next.Lexeme)),
				)
			}

			items = append(items, infixItem{op: string(next.Lexeme)}, infixItem{expr: rhs})
		}

		return resolveInfix(items, p.fixity), nil
	}
}

func (p *Parser[S]) appexp(
	lex tokens.Lexer[Kind, Lexeme, S],
	steps []parsec.Combinator[rune, strings.Position, Expr, S],
	rec parsec.Combinator[rune, strings.Position, Expr, S],
) parsec.Combinator[rune, strings.Position, Expr, S] {
	first := p.atom(lex, steps, true, rec)
	args := parsec.Many(0, parsec.Try(p.atom(lex, steps, false, rec)))

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		e, err := first(state, buf)
		if err != nil {
			return nil, err
		}

		rest, _ := args(state, buf)

		for _, arg := range rest {
			e = App{Fn: e, Arg: arg}
		}

		return e, nil
	}
}

func (p *Parser[S]) atom(
	lex tokens.Lexer[Kind, Lexeme, S],
	steps []parsec.Combinator[rune, strings.Position, Expr, S],
	operandStart bool,
	rec parsec.Combinator[rune, strings.Position, Expr, S],
) parsec.Combinator[rune, strings.Position, Expr, S] {
	primary := p.primary(lex, steps, operandStart, rec)
	postfix := parsec.Many(0, parsec.Try(tokens.Satisfy(lex, "postfix operator", func(t Token) bool {
		_, isPostfix := p.postfix[string(t.Lexeme)]

		return t.Kind == KindOperator && isPostfix
	})))

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		e, err := primary(state, buf)
		if err != nil {
			return nil, err
		}

		ops, _ := postfix(state, buf)

		for _, t := range ops {
			e = Postfix{Op: string(t.Lexeme), E: e}
		}

		return e, nil
	}
}

func (p *Parser[S]) primary(
	lex tokens.Lexer[Kind, Lexeme, S],
	steps []parsec.Combinator[rune, strings.Position, Expr, S],
	operandStart bool,
	rec parsec.Combinator[rune, strings.Position, Expr, S],
) parsec.Combinator[rune, strings.Position, Expr, S] {
	intLit := parsec.Cast(tokens.OfKind(lex, KindInt), func(t Token) (Expr, error) {
		value, _ := strconv.ParseInt(string(t.Lexeme), 10, 64)

		return IntLit{Text: string(t.Lexeme), Value: value}, nil
	})

	ident := parsec.Cast(tokens.OfKind(lex, KindIdent), func(t Token) (Expr, error) {
		return Ident{Name: string(t.Lexeme)}, nil
	})

	prefixOp := tokens.Satisfy(lex, "prefix operator", func(t Token) bool {
		_, isPrefix := p.prefix[string(t.Lexeme)]
		if !isPrefix || t.Kind != KindOperator {
			return false
		}

		// a - b is infix, not a applied to -b: inside an argument a
		// prefix operator must not carry fixity status
		if !operandStart {
			_, hasFixity := p.fixity[string(t.Lexeme)]

			return !hasFixity
		}

		return true
	})

	open := tokens.Exact(lex, KindSymbol, "(")
	closeParen := tokens.Exact(lex, KindSymbol, ")")

	// rec is the full expression, the expression forms included: a
	// parenthesized expression is again a full one
	paren := func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		if _, err := open(state, buf); err != nil {
			return nil, err
		}

		e, err := rec(state, buf)
		if err != nil {
			return nil, err
		}

		if _, err := closeParen(state, buf); err != nil {
			return nil, err
		}

		return e, nil
	}

	prefixExpr := func(state S, buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		t, err := prefixOp(state, buf)
		if err != nil {
			return nil, err
		}

		sub, serr := p.atom(lex, steps, true, rec)(state, buf)
		if serr != nil {
			return nil, serr
		}

		return Unary{Op: string(t.Lexeme), Sub: sub}, nil
	}

	alts := make([]parsec.Combinator[rune, strings.Position, Expr, S], 0, len(steps)+4)

	for _, step := range steps {
		alts = append(alts, parsec.Try(step))
	}

	alts = append(alts,
		parsec.Try(prefixExpr),
		parsec.Try(intLit),
		parsec.Try(ident),
		parsec.Try(paren),
	)

	return parsec.Choice("expected operand", alts...)
}

type infixItem struct {
	op   string
	expr Expr
}

type ternaryDef struct {
	colon string
	prec  int
}

func resolveInfix(items []infixItem, fix map[string]Fixity) Expr {
	reduced := reduceInfix(items, fix, Fixity{Prec: math.MinInt})

	return reduced[0].expr
}

// reduceInfix resolves the operators tighter than cur wherever they
// sit in the chain and leaves the operators of the precedence of cur
// and below as items between the reduced expressions: x = a + 2 ?
// b : c becomes x = (+ a 2) at the boundary of the conditional, its
// condition the last expression of the reduced chain.
func reduceInfix(items []infixItem, fix map[string]Fixity, cur Fixity) []infixItem {
	reduced := make([]infixItem, 0, len(items))

	var vals []Expr
	var ops []string

	apply := func() {
		r := vals[len(vals)-1]
		l := vals[len(vals)-2]
		vals = vals[:len(vals)-2]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		vals = append(vals, Binary{Op: op, L: l, R: r})
	}

	flush := func() {
		for len(ops) > 0 {
			apply()
		}

		for _, v := range vals {
			reduced = append(reduced, infixItem{expr: v})
		}

		vals = vals[:0]
	}

	for _, it := range items {
		if it.op == "" {
			vals = append(vals, it.expr)
			continue
		}

		// an operator of the precedence of cur or below keeps its place
		// in the chain — the = of x = a + 2 ? b : c waits for the whole
		// conditional — while the tighter ones around it still reduce
		if fix[it.op].Prec <= cur.Prec {
			flush()
			reduced = append(reduced, it)

			continue
		}

		itFix := fix[it.op]

		for len(ops) > 0 {
			top := fix[ops[len(ops)-1]]
			if top.Prec > itFix.Prec || top.Prec == itFix.Prec && !itFix.Right {
				apply()
			} else {
				break
			}
		}

		ops = append(ops, it.op)
	}

	flush()

	return reduced
}
