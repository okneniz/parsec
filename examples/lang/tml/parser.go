package tml

import (
	"errors"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/lang"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

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
	expr parsec.Combinator[rune, Position, Expr]
}

func newParser() *parser {
	p := &parser{}

	ops := lang.NewParser().
		Infix(7, "*", "/").
		Infix(6, "+", "-").
		Infix(4, "=", "<>", ">", ">=", "<", "<=").
		InfixRight(-1, "andalso").
		InfixRight(-2, "orelse").
		ExprForm(p.ifForm, p.fnForm).
		PrefixForm(p.letForm, p.negForm, p.parenForm)

	p.expr = ops.Expr(lexer, boolOperand("true", true), boolOperand("false", false))

	return p
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

	return func(buf parsec.Buffer[rune, Position]) (Decl, parsec.Error[Position]) {
		pat, perr := atpat(buf)
		if perr != nil {
			return nil, perr
		}

		if _, err := eq(buf); err != nil {
			return nil, err
		}

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		e, eerr := p.expr(buf)
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

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		body, berr := p.expr(buf)
		if berr != nil {
			return nil, berr
		}
		decl.Body = body

		return decl, nil
	}
}

// ---------------------------------------------------------------------------
// Expression forms

// ifForm builds the if-expression form, if exp then exp else exp. It
// is a form of the expression level — an if is not an operand, and
// every branch is again a full expression.
func (p *parser) ifForm(l lang.Ladder) parsec.Combinator[rune, Position, Expr] {
	ifKw := tokens.Exact(l.Lex, KindKeyword, "if")
	thenKw := tokens.Exact(l.Lex, KindKeyword, "then")
	elseKw := tokens.Exact(l.Lex, KindKeyword, "else")

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := ifKw(buf); err != nil {
			return nil, err
		}

		cond, cerr := l.Expr(buf)
		if cerr != nil {
			return nil, cerr
		}

		if _, err := thenKw(buf); err != nil {
			return nil, err
		}

		then, terr := l.Expr(buf)
		if terr != nil {
			return nil, terr
		}

		if _, err := elseKw(buf); err != nil {
			return nil, err
		}

		els, eerr := l.Expr(buf)
		if eerr != nil {
			return nil, eerr
		}

		return If{Cond: cond, Then: then, Else: els}, nil
	}
}

// fnForm builds the function form, fn pat => exp. It is a form of
// the expression level, like if: a lambda is not an operand.
func (p *parser) fnForm(l lang.Ladder) parsec.Combinator[rune, Position, Expr] {
	fnKw := tokens.Exact(l.Lex, KindKeyword, "fn")
	arrow := tokens.Exact(l.Lex, KindOperator, "=>")
	atpat := p.atpat(l.Lex)

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := fnKw(buf); err != nil {
			return nil, err
		}

		arg, aerr := atpat(buf)
		if aerr != nil {
			return nil, aerr
		}

		if _, err := arrow(buf); err != nil {
			return nil, err
		}

		body, berr := l.Expr(buf)
		if berr != nil {
			return nil, berr
		}

		return Fn{Arg: arg, Body: body}, nil
	}
}

// letForm builds the let-expression form, decl+ in exp end. It is a
// form of the operand level — a let is an atom, f let ... end applies
// f — and its body is a full expression.
func (p *parser) letForm(l lang.Ladder) parsec.Combinator[rune, Position, Expr] {
	letKw := tokens.Exact(l.Lex, KindKeyword, "let")
	inKw := tokens.Exact(l.Lex, KindKeyword, "in")
	endKw := tokens.Exact(l.Lex, KindKeyword, "end")
	decs := p.decList(l.Lex)

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := letKw(buf); err != nil {
			return nil, err
		}

		decls, derr := decs(buf)
		if derr != nil {
			return nil, derr
		}

		if len(decls) == 0 {
			return nil, parsec.NewParseError(buf.Position(), "expected declaration in let")
		}

		if _, err := inKw(buf); err != nil {
			return nil, err
		}

		body, berr := l.Expr(buf)
		if berr != nil {
			return nil, berr
		}

		if _, err := endKw(buf); err != nil {
			return nil, err
		}

		return Let{Decls: decls, Body: body}, nil
	}
}

// negForm builds the negated literal ~3, with the negation glued to
// the literal so that f ~3 is f applied to ~3, not f ~ applied to 3.
// A lone ~ is the negation function itself.
func (p *parser) negForm(l lang.Ladder) parsec.Combinator[rune, Position, Expr] {
	tilde := tokens.Exact(l.Lex, KindOperator, "~")
	intLit := parsec.Try(tokens.OfKind(l.Lex, KindInt))

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := tilde(buf); err != nil {
			return nil, err
		}

		lit, lerr := intLit(buf)
		if lerr != nil {
			return Ident{Name: "~"}, nil
		}

		value, _ := strconv.ParseInt(string(lit.Lexeme), 10, 64)

		return App{Fn: Ident{Name: "~"}, Arg: IntLit{Text: string(lit.Lexeme), Value: value}}, nil
	}
}

// parenForm builds the grouping form (e) and the tuple form
// (e1, ..., en); every element is again a full expression.
func (p *parser) parenForm(l lang.Ladder) parsec.Combinator[rune, Position, Expr] {
	open := tokens.Exact(l.Lex, KindSymbol, "(")
	closeParen := tokens.Exact(l.Lex, KindSymbol, ")")
	comma := parsec.Try(tokens.Exact(l.Lex, KindSymbol, ","))

	return func(buf parsec.Buffer[rune, Position]) (Expr, parsec.Error[Position]) {
		if _, err := open(buf); err != nil {
			return nil, err
		}

		first, ferr := l.Expr(buf)
		if ferr != nil {
			return nil, ferr
		}

		items := []Expr{first}

		for {
			if _, cerr := comma(buf); cerr != nil {
				break
			}

			e, eerr := l.Expr(buf)
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

// boolOperand turns a boolean literal into an extra operand of the
// ladder: a keyword with a fixed meaning needs no recursion.
func boolOperand(lexeme Lexeme, value bool) lang.Operand {
	return func(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, Position, Expr] {
		return parsec.Cast(
			tokens.Exact(lex, KindKeyword, lexeme),
			func(Token) (Expr, error) {
				return BoolLit{Value: value}, nil
			},
		)
	}
}

// ---------------------------------------------------------------------------
// Patterns

// atpat parses one atomic pattern. Tiny ml has no constructors, so
// there are no infix constructor patterns and no layered patterns.
// The dispatch mirrors the operand level of the ladder: a Choice of
// Try-wrapped alternatives.
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
