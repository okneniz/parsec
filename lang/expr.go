package lang

import (
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

type Fixity struct {
	Prec  int
	Right bool
}

type Parser struct {
	fixity      map[string]Fixity
	prefix      map[string]struct{}
	postfix     map[string]struct{}
	exprForms   []Form
	prefixForms []Form
}

// NewParser creates a parser with an empty operator table.
func NewParser() *Parser {
	return &Parser{
		fixity:  map[string]Fixity{},
		prefix:  map[string]struct{}{},
		postfix: map[string]struct{}{},
	}
}

// Infix declares left-associative operators of the given precedence;
// higher binds tighter.
func (p *Parser) Infix(prec int, ops ...string) *Parser {
	for _, op := range ops {
		p.fixity[op] = Fixity{Prec: prec}
	}
	return p
}

// InfixRight declares right-associative operators of the given
// precedence, as exponentiation usually is.
func (p *Parser) InfixRight(prec int, ops ...string) *Parser {
	for _, op := range ops {
		p.fixity[op] = Fixity{Prec: prec, Right: true}
	}
	return p
}

// Prefix declares prefix operators. A prefix operator applies to the
// next operand only, so -a + b is (+ (- a) b) and -(a + b) needs
// parentheses.
func (p *Parser) Prefix(ops ...string) *Parser {
	for _, op := range ops {
		p.prefix[op] = struct{}{}
	}
	return p
}

// Postfix declares postfix operators. A postfix operator applies to
// the complete operand before it and chains right, so a! + b is
// (+ (! a) b) and a!! is (! (! a)). An operator must not be both
// postfix and infix: the postfix loop would consume it first.
func (p *Parser) Postfix(ops ...string) *Parser {
	for _, op := range ops {
		p.postfix[op] = struct{}{}
	}
	return p
}

// ExprForm declares forms of the expression level, tried before the
// infix chain: the if e then e else e and fn p => e of sml live
// here, above every operator. An expression form is not an operand —
// 1 + if e then e else e does not parse — and it consumes its whole
// rest itself, subexpressions included, through the ladder.
func (p *Parser) ExprForm(forms ...Form) *Parser {
	p.exprForms = append(p.exprForms, forms...)

	return p
}

// PrefixForm declares forms of the operand level, tried among the
// operands before the built-in ones: a sizeof e reads its operand
// through the ladder's Atom and so binds tighter than any infix
// operator, a cast (t) e likewise; a form that owns its
// subexpressions wholly — a let, a parenthesized tuple — parses them
// through the ladder's Expr.
func (p *Parser) PrefixForm(forms ...Form) *Parser {
	p.prefixForms = append(p.prefixForms, forms...)

	return p
}

// Ladder hands a form the recursive entry points of the expression
// machinery, so the form parses its subexpressions with the same
// parser it is part of. Expr is the full expression, forms included;
// Atom is one operand, with its prefix and postfix operators — the
// level a prefix operator applies to. Lex is the lexer the whole
// ladder reads through.
type Ladder struct {
	Lex  tokens.Lexer[Kind, Lexeme]
	Expr parsec.Combinator[rune, strings.Position, Expr]
	Atom parsec.Combinator[rune, strings.Position, Expr]
}

// A Form builds the combinator of a syntactic form that does not
// reduce to an operator lexeme: a keyword-led form such as sizeof e
// or if e then e else e, or a form led by a punctuation token, a
// cast (t) e, a parenthesized tuple. The forms of a Parser are built
// once, when its Expr runs; everything a form needs besides the
// ladder it builds above the returned closure.
type Form func(Ladder) parsec.Combinator[rune, strings.Position, Expr]

type Operand func(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, strings.Position, Expr]

func (p *Parser) Expr(
	lex tokens.Lexer[Kind, Lexeme],
	extras ...Operand,
) parsec.Combinator[rune, strings.Position, Expr] {
	steps := make([]parsec.Combinator[rune, strings.Position, Expr], len(extras), len(extras)+len(p.prefixForms))

	for i, extra := range extras {
		steps[i] = extra(lex)
	}

	// the forms are part of the ladder they recurse into, so the
	// ladder hands out closures that read full and atom at call
	// time, after the whole ladder is built
	ladder := Ladder{Lex: lex}

	var full parsec.Combinator[rune, strings.Position, Expr]
	var atom parsec.Combinator[rune, strings.Position, Expr]

	ladder.Expr = func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		return full(buf)
	}

	ladder.Atom = func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		return atom(buf)
	}

	for _, form := range p.prefixForms {
		steps = append(steps, form(ladder))
	}

	full = p.expr(lex, steps, ladder.Expr)

	if len(p.exprForms) > 0 {
		alts := make([]parsec.Combinator[rune, strings.Position, Expr], 0, len(p.exprForms)+1)

		for _, form := range p.exprForms {
			alts = append(alts, parsec.Try(form(ladder)))
		}

		full = parsec.Choice("expected expression", append(alts, full)...)
	}

	atom = p.atom(lex, steps, true, ladder.Expr)

	return full
}

func (p *Parser) expr(
	lex tokens.Lexer[Kind, Lexeme],
	steps []parsec.Combinator[rune, strings.Position, Expr],
	rec parsec.Combinator[rune, strings.Position, Expr],
) parsec.Combinator[rune, strings.Position, Expr] {
	// an operator is a lexeme with fixity status, whatever kind the
	// lexer gave its token: the reserved words andalso and orelse of
	// sml are infix operators of their own precedence level
	op := parsec.Try(tokens.Satisfy(lex, "operator", func(t Token) bool {
		_, hasFixity := p.fixity[string(t.Lexeme)]

		return hasFixity
	}))

	return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		first, err := p.appexp(lex, steps, rec)(buf)
		if err != nil {
			return nil, err
		}

		items := []infixItem{{expr: first}}

		for {
			next, oerr := op(buf)
			if oerr != nil {
				break
			}

			rhs, rerr := p.appexp(lex, steps, rec)(buf)
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

func (p *Parser) appexp(
	lex tokens.Lexer[Kind, Lexeme],
	steps []parsec.Combinator[rune, strings.Position, Expr],
	rec parsec.Combinator[rune, strings.Position, Expr],
) parsec.Combinator[rune, strings.Position, Expr] {
	first := p.atom(lex, steps, true, rec)
	args := parsec.Many(0, parsec.Try(p.atom(lex, steps, false, rec)))

	return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
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

func (p *Parser) atom(
	lex tokens.Lexer[Kind, Lexeme],
	steps []parsec.Combinator[rune, strings.Position, Expr],
	operandStart bool,
	rec parsec.Combinator[rune, strings.Position, Expr],
) parsec.Combinator[rune, strings.Position, Expr] {
	primary := p.primary(lex, steps, operandStart, rec)
	postfix := parsec.Many(0, parsec.Try(tokens.Satisfy(lex, "postfix operator", func(t Token) bool {
		_, isPostfix := p.postfix[string(t.Lexeme)]

		return t.Kind == KindOperator && isPostfix
	})))

	return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		e, err := primary(buf)
		if err != nil {
			return nil, err
		}

		ops, _ := postfix(buf)

		for _, t := range ops {
			e = Postfix{Op: string(t.Lexeme), E: e}
		}

		return e, nil
	}
}

func (p *Parser) primary(
	lex tokens.Lexer[Kind, Lexeme],
	steps []parsec.Combinator[rune, strings.Position, Expr],
	operandStart bool,
	rec parsec.Combinator[rune, strings.Position, Expr],
) parsec.Combinator[rune, strings.Position, Expr] {
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
	paren := func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		if _, err := open(buf); err != nil {
			return nil, err
		}

		e, err := rec(buf)
		if err != nil {
			return nil, err
		}

		if _, err := closeParen(buf); err != nil {
			return nil, err
		}

		return e, nil
	}

	prefixExpr := func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		t, err := prefixOp(buf)
		if err != nil {
			return nil, err
		}

		sub, serr := p.atom(lex, steps, true, rec)(buf)
		if serr != nil {
			return nil, serr
		}

		return Unary{Op: string(t.Lexeme), Sub: sub}, nil
	}

	alts := make([]parsec.Combinator[rune, strings.Position, Expr], 0, len(steps)+4)

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

func resolveInfix(items []infixItem, fix map[string]Fixity) Expr {
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

	for _, it := range items {
		if it.op == "" {
			vals = append(vals, it.expr)
			continue
		}

		cur := fix[it.op]
		for len(ops) > 0 {
			top := fix[ops[len(ops)-1]]
			if top.Prec > cur.Prec || top.Prec == cur.Prec && !cur.Right {
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
