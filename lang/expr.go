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
	fixity  map[string]Fixity
	prefix  map[string]struct{}
	postfix map[string]struct{}
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

type Operand func(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, strings.Position, Expr]

func (p *Parser) Expr(
	lex tokens.Lexer[Kind, Lexeme],
	extras ...Operand,
) parsec.Combinator[rune, strings.Position, Expr] {
	steps := make([]parsec.Combinator[rune, strings.Position, Expr], len(extras))

	for i, extra := range extras {
		steps[i] = extra(lex)
	}

	return p.expr(lex, steps)
}

func (p *Parser) expr(
	lex tokens.Lexer[Kind, Lexeme],
	steps []parsec.Combinator[rune, strings.Position, Expr],
) parsec.Combinator[rune, strings.Position, Expr] {
	op := parsec.Try(tokens.Satisfy(lex, "operator", func(t Token) bool {
		_, hasFixity := p.fixity[string(t.Lexeme)]

		return t.Kind == KindOperator && hasFixity
	}))

	return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		first, err := p.appexp(lex, steps)(buf)
		if err != nil {
			return nil, err
		}

		items := []infixItem{{expr: first}}

		for {
			next, oerr := op(buf)
			if oerr != nil {
				break
			}

			rhs, rerr := p.appexp(lex, steps)(buf)
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
) parsec.Combinator[rune, strings.Position, Expr] {
	first := p.atom(lex, steps, true)
	args := parsec.Many(0, parsec.Try(p.atom(lex, steps, false)))

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
) parsec.Combinator[rune, strings.Position, Expr] {
	primary := p.primary(lex, steps, operandStart)
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

	paren := func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		if _, err := open(buf); err != nil {
			return nil, err
		}

		e, err := p.expr(lex, steps)(buf)
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

		sub, serr := p.atom(lex, steps, true)(buf)
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
