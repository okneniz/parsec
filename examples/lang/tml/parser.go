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
// ladder from a form, into atpat from parseTuplePattern — come with the
// ladder or stay in the body.
type parser struct {
	parseExpression Parser[Expr]
	lexer           Lexer
}

type Lexer = tokens.Lexer[Kind, Lexeme, *state]
type Parser[T any] = parsec.Combinator[rune, Position, T, *state]
type Error = parsec.Error[Position]

func newParser() *parser {
	p := &parser{
		lexer: lexer,
	}

	ops := lang.NewParser[*state]().
		Infix(7, "*", "/").
		Infix(6, "+", "-").
		Infix(4, "=", "<>", ">", ">=", "<", "<=").
		InfixRight(-1, "andalso").
		InfixRight(-2, "orelse").
		ExprForm(p.makeConditionParser, p.makeFnFormParser).
		PrefixForm(p.makeLetFormParser, p.makeNegatedFormParser, p.parenForm)

	p.parseExpression = ops.Expr(
		lexer,
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
	parseAST := newParser().makeDeclarationsParser()

	ast, err := parseAST(st, buf)
	if err != nil {
		return nil, err
	}

	// TODO : looks like bad design
	fixup(ast)

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

	return ast, nil
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

// makeDeclarationsParser parses declarations until the next token cannot start one.
// The caller checks its own end condition: the end of input, or the
// in keyword of a let. The leading token selects the form from a
// MapAs table keyed by the whole token and is discarded; a new
// declaration form is a new table key. A declaration that starts but
// does not parse stops the list as well — the caller's end check
// reports it.
func (p *parser) makeDeclarationsParser() Parser[[]Decl] {
	parseValDeclaration := p.makeValDeclarationParser()
	parseFuncDeclaration := p.makeFuncDeclarationParser()

	forms := map[Token]Parser[Decl]{
		{Kind: KindKeyword, Lexeme: "val"}: parseValDeclaration,
		{Kind: KindKeyword, Lexeme: "fun"}: parseFuncDeclaration,
	}

	parseDeclaration := parsec.MapAs(
		"expected declaration",
		forms,
		// TODO : looks strange
		tokens.Satisfy(p.lexer, "expected declaration", parsec.Anything[Token]),
	)

	parseDeclaration = parsec.Try(parseDeclaration)

	return parsec.Many(0, parseDeclaration)
}

// makeValDeclarationParser parses the rest of a val declaration, pat = exp; the val
// keyword is already consumed. The expression is read in the scope
// the declaration appears in; its pattern binds over everything that
// follows, so the binders are pushed only after it.
func (p *parser) makeValDeclarationParser() Parser[Decl] {
	parseEq := tokens.Exact(p.lexer, KindOperator, "=")
	parseAtomicPattern := p.makeAtomicPatternParser()

	return func(st *state, buf parsec.Buffer[rune, Position]) (Decl, Error) {
		pat, perr := parseAtomicPattern(st, buf)
		if perr != nil {
			return nil, perr
		}

		if _, err := parseEq(st, buf); err != nil {
			return nil, err
		}

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		e, eerr := p.parseExpression(st, buf)
		if eerr != nil {
			return nil, eerr
		}

		st.push(patVars(pat)...)

		return ValDecl{Binds: []ValBind{{Pat: pat, E: e}}}, nil
	}
}

// makeFuncDeclarationParser parses the rest of a recursive function declaration, fun
// f x1 ... xn = exp with atomic argument patterns; the fun keyword
// is already consumed. The name and the arguments bind inside the
// body — recursion — and, like every declaration binder, over
// everything that follows.
func (p *parser) makeFuncDeclarationParser() Parser[Decl] {
	parseEq := tokens.Exact(p.lexer, KindOperator, "=")
	parseIdent := tokens.OfKind(p.lexer, KindIdent)
	parseArg := p.makeAtomicPatternParser()
	parseArg = parsec.Try(parseArg)
	parseArgs := parsec.Many(0, parseArg)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Decl, Error) {
		name, nerr := parseIdent(st, buf)
		if nerr != nil {
			return nil, nerr
		}

		args, err := parseArgs(st, buf)
		if err != nil {
			return nil, err
		}

		if _, err := parseEq(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()

		st.push(string(name.Lexeme))
		st.push(patVarsAll(args)...)

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		body, berr := p.parseExpression(st, buf)
		if berr != nil {
			st.popTo(base)

			return nil, berr
		}

		// the arguments scope over the body only; the name binds
		// over the rest of the program
		st.popTo(base + 1)

		decl := FunDecl{
			Name: string(name.Lexeme),
			Args: args,
			Body: body,
		}

		return decl, nil
	}
}

// ---------------------------------------------------------------------------
// Expression forms

// makeConditionParser builds the if-expression form, if exp then exp else exp. It
// is a form of the expression level — an if is not an operand, and
// every branch is again a full expression.
func (p *parser) makeConditionParser(ladder lang.Ladder[*state]) Parser[Expr] {
	parseIf := tokens.Exact(ladder.Lex, KindKeyword, "if")
	parseThen := tokens.Exact(ladder.Lex, KindKeyword, "then")
	parseElse := tokens.Exact(ladder.Lex, KindKeyword, "else")

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseIf(st, buf); err != nil {
			return nil, err
		}

		cond, cerr := ladder.Expr(st, buf)
		if cerr != nil {
			return nil, cerr
		}

		if _, err := parseThen(st, buf); err != nil {
			return nil, err
		}

		then, terr := ladder.Expr(st, buf)
		if terr != nil {
			return nil, terr
		}

		if _, err := parseElse(st, buf); err != nil {
			return nil, err
		}

		els, eerr := ladder.Expr(st, buf)
		if eerr != nil {
			return nil, eerr
		}

		return If{
			Cond: cond,
			Then: then,
			Else: els,
		}, nil
	}
}

// makeFnFormParser builds the function form, fn pat => exp. It is a form of
// the expression level, like if: a lambda is not an operand. The
// pattern binds over the body — pushed before it and popped after,
// on the failure paths too, because a retried alternative must see
// the stack as it was.
func (p *parser) makeFnFormParser(ladder lang.Ladder[*state]) Parser[Expr] {
	parseFn := tokens.Exact(ladder.Lex, KindKeyword, "fn")
	parseArrow := tokens.Exact(ladder.Lex, KindOperator, "=>")
	parseAtomicPattern := p.makeAtomicPatternParser()

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseFn(st, buf); err != nil {
			return nil, err
		}

		arg, aerr := parseAtomicPattern(st, buf)
		if aerr != nil {
			return nil, aerr
		}

		if _, err := parseArrow(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()

		st.push(patVars(arg)...)
		defer st.popTo(base)

		body, berr := ladder.Expr(st, buf)
		if berr != nil {
			return nil, berr
		}

		return Fn{
			Arg:  arg,
			Body: body,
		}, nil
	}
}

// makeLetFormParser builds the let-expression form, decl+ in exp end. It is
// a form of the operand level — a let is an atom, f let ... end applies
// f — and its body is a full expression. The declarations bind
// sequentially over the body; the whole block pops its binders on
// the way out, on the failure paths too.
func (p *parser) makeLetFormParser(ladder lang.Ladder[*state]) Parser[Expr] {
	parseLet := tokens.Exact(ladder.Lex, KindKeyword, "let")
	parseIn := tokens.Exact(ladder.Lex, KindKeyword, "in")
	parseEnd := tokens.Exact(ladder.Lex, KindKeyword, "end")
	parseDeclarations := p.makeDeclarationsParser()

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseLet(st, buf); err != nil {
			return nil, err
		}

		base := st.depth()
		defer st.popTo(base)

		decls, derr := parseDeclarations(st, buf)
		if derr != nil {
			return nil, derr
		}

		if len(decls) == 0 {
			return nil, parsec.NewParseError(buf.Position(), "expected declaration in let")
		}

		if _, err := parseIn(st, buf); err != nil {
			return nil, err
		}

		body, berr := ladder.Expr(st, buf)
		if berr != nil {
			return nil, berr
		}

		if _, err := parseEnd(st, buf); err != nil {
			return nil, err
		}

		return Let{
			Decls: decls,
			Body:  body,
		}, nil
	}
}

// makeNegatedFormParser builds the negated literal ~3, with the negation glued to
// the literal so that f ~3 is f applied to ~3, not f ~ applied to 3.
// A lone ~ is the negation function itself — the builtin binder of
// the initial stack.
func (p *parser) makeNegatedFormParser(l lang.Ladder[*state]) Parser[Expr] {
	// TODO : add error message to Exact

	parseTilde := tokens.Exact(l.Lex, KindOperator, "~")
	parseLiteral := parsec.Try(tokens.OfKind(l.Lex, KindInt))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseTilde(st, buf); err != nil {
			return nil, err
		}

		pos := buf.Position()
		neg := st.builtinVar()

		lit, lerr := parseLiteral(st, buf)
		if lerr != nil {
			return neg, nil
		}

		value, castErr := strconv.ParseInt(string(lit.Lexeme), 10, 64)
		if castErr != nil {
			return nil, parsec.NewParseError(pos, castErr.Error())
		}

		return App{
			Fn: neg,
			Arg: IntLit{
				Text:  string(lit.Lexeme),
				Value: value,
			},
		}, nil
	}
}

// parenForm builds the grouping form (e) and the tuple form
// (e1, ..., en); every element is again a full expression.
func (p *parser) parenForm(ladder lang.Ladder[*state]) Parser[Expr] {
	parseLeftParen := tokens.Exact(ladder.Lex, KindSymbol, "(")
	parseRightParen := tokens.Exact(ladder.Lex, KindSymbol, ")")
	parseComma := parsec.Try(tokens.Exact(ladder.Lex, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseLeftParen(st, buf); err != nil {
			return nil, err
		}

		first, ferr := ladder.Expr(st, buf)
		if ferr != nil {
			return nil, ferr
		}

		items := []Expr{first}

		// TODO : why not sepEndBy?
		for {
			if _, cerr := parseComma(st, buf); cerr != nil {
				break
			}

			e, eerr := ladder.Expr(st, buf)
			if eerr != nil {
				return nil, eerr
			}

			items = append(items, e)
		}

		if _, err := parseRightParen(st, buf); err != nil {
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
	return func(lex Lexer) Parser[Expr] {
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
func varOperand(lex Lexer) Parser[Expr] {
	ident := tokens.OfKind(lex, KindIdent)

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
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

// makeAtomicPatternParser parses one atomic pattern. Tiny ml has no constructors, so
// there are no infix constructor patterns and no layered patterns.
// The dispatch mirrors the operand level of the ladder: a Choice of
// Try-wrapped alternatives.
func (p *parser) makeAtomicPatternParser() parsec.Combinator[rune, Position, Pat, *state] {
	parseWildcard := tokens.Exact(p.lexer, KindSymbol, "_")

	toWildcard := parsec.Cast(
		parseWildcard,
		func(Token) (Pat, error) {
			return WildcardPat{}, nil
		},
	)

	parseConstPattern := tokens.Satisfy(
		p.lexer,
		"expected literal pattern",
		isLiteralToken,
	)

	toConstPattern := parsec.Cast(
		parseConstPattern,
		func(t Token) (Pat, error) {
			return ConstPat{Text: string(t.Lexeme)}, nil
		},
	)

	parseVarPattern := tokens.OfKind(p.lexer, KindIdent)
	toVarPattern := parsec.Cast(
		parseVarPattern,
		func(t Token) (Pat, error) {
			return VarPat{Name: string(t.Lexeme)}, nil
		},
	)

	parseNegativeLiteralPattern := p.makeNegativeLiteralPattern()
	parseTuplePattern := p.makeTuplePatternParser()

	return parsec.Choice("expected pattern",
		parsec.Try(toWildcard),
		parsec.Try(toConstPattern),
		parsec.Try(parseNegativeLiteralPattern),
		parsec.Try(toVarPattern),
		parsec.Try(parseTuplePattern),
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

// makeNegativeLiteralPattern parses the negative literal pattern ~1; a lone ~ is a
// variable pattern.
func (p *parser) makeNegativeLiteralPattern() parsec.Combinator[rune, Position, Pat, *state] {
	parseTilde := tokens.Exact(p.lexer, KindOperator, "~")
	parseIntLiteral := parsec.Try(tokens.OfKind(p.lexer, KindInt))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Pat, Error) {
		if _, err := parseTilde(st, buf); err != nil {
			return nil, err
		}

		lit, lerr := parseIntLiteral(st, buf)
		if lerr != nil {
			return VarPat{Name: "~"}, nil
		}

		return ConstPat{Text: "~" + string(lit.Lexeme)}, nil
	}
}

// makeTuplePatternParser parses (p) and tuple patterns (p1, ..., pn).
func (p *parser) makeTuplePatternParser() parsec.Combinator[rune, Position, Pat, *state] {
	parseLeftParen := tokens.Exact(p.lexer, KindSymbol, "(")
	parseRightParen := tokens.Exact(p.lexer, KindSymbol, ")")
	parseComma := parsec.Try(tokens.Exact(p.lexer, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Pat, Error) {
		if _, err := parseLeftParen(st, buf); err != nil {
			return nil, err
		}

		// TODO : why not use scope / closure?
		//
		// atpat's Choice builds this form, so the reference stays
		// in the body; one build serves every element
		parseAtomicPattern := p.makeAtomicPatternParser()

		var items []Pat

		// TODO : why not just sepEndBy?
		for {
			pat, perr := parseAtomicPattern(st, buf)
			if perr != nil {
				return nil, perr
			}

			items = append(items, pat)

			if _, cerr := parseComma(st, buf); cerr != nil {
				break
			}
		}

		if _, err := parseRightParen(st, buf); err != nil {
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
