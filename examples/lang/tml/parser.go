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

type (
	Lexer         = tokens.Lexer[Kind, Lexeme, *state]
	Parser[T any] = parsec.Combinator[rune, Position, T, *state]
	Error         = parsec.Error[Position]
)

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
// val and fun declarations, resolved to de Bruijn indices. The top
// level binds in two passes: the first collects the names of the
// functions, the second — with the names already on the binder stack
// — resolves every reference to them, so the functions see each
// other regardless of declaration order.
func Parse(src string) ([]Decl, error) {
	first, err := parseProgram(src, newState())
	if err != nil {
		return nil, err
	}

	funs := funNames(first)

	// no functions to pre-bind: the second pass would start from the
	// same state as the first and parse the same program
	if len(funs) == 0 {
		return first, nil
	}

	st := newState()
	st.push(funs...)

	return parseProgram(src, st)
}

// funNames lists the names of the top-level functions in declaration
// order: the row the second parse pass pre-binds on the binder stack.
func funNames(decls []Decl) []string {
	funs := []string{}

	for _, d := range decls {
		if fun, ok := d.(FunDecl); ok {
			funs = append(funs, fun.Name)
		}
	}

	return funs
}

// parseProgram runs the declarations parser over a fresh buffer with
// the given state and checks that nothing is left of the input. The
// funs bind through the pre-bound row of the state.
func parseProgram(src string, st *state) ([]Decl, error) {
	buf := strings.Buffer([]rune(src))
	p := newParser()
	parseAST := p.makeDeclarationsParser(p.makePreboundFuncDeclarationParser())

	ast, err := parseAST(st, buf)
	if err != nil {
		return nil, err
	}

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

// ---------------------------------------------------------------------------
// Declarations

// makeDeclarationsParser parses declarations until the next token
// cannot start one. The caller checks its own end condition: the end
// of input, or the in keyword of a let. The leading token selects
// the form from a MapAs table keyed by the whole token and is
// discarded; a new declaration form is a new table key. A declaration
// that starts but does not parse stops the list as well — the
// caller's end check reports it. The fun form arrives as an
// argument: a top-level fun binds through its pre-bound slot, a let
// fun pushes its own name.
func (p *parser) makeDeclarationsParser(parseFuncDeclaration Parser[Decl]) Parser[[]Decl] {
	parseValDeclaration := p.makeValDeclarationParser()

	forms := map[Token]Parser[Decl]{
		{Kind: KindKeyword, Lexeme: "val"}: parseValDeclaration,
		{Kind: KindKeyword, Lexeme: "fun"}: parseFuncDeclaration,
	}

	parseDeclaration := parsec.MapAs(
		"expected declaration",
		forms,
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

// makeFuncDeclarationParser parses the rest of a recursive function
// declaration, fun f x1 ... xn = exp with atomic argument patterns;
// the fun keyword is already consumed. The name and the arguments
// bind inside the body — recursion — and the name, like every
// declaration binder, over everything that follows: this is the fun
// of a let body.
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

// makePreboundFuncDeclarationParser is the top-level fun. Its slot
// sits in the pre-bound fun row the second parse pass starts with, so
// the shadow it pushes for the body — and the arguments with it —
// pops away completely.
func (p *parser) makePreboundFuncDeclarationParser() Parser[Decl] {
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
		defer st.popTo(base)

		// the ladder is assigned after the declarations are built —
		// the field is read here, at call time
		body, berr := p.parseExpression(st, buf)
		if berr != nil {
			return nil, berr
		}

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
	parseDeclarations := p.makeDeclarationsParser(p.makeFuncDeclarationParser())

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
	parseEmptyTuple := parsec.Try(parseRightParen)
	parseComma := parsec.Try(tokens.Exact(ladder.Lex, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Expr, Error) {
		if _, err := parseLeftParen(st, buf); err != nil {
			return nil, err
		}

		// the unit () is the empty tuple
		if _, err := parseEmptyTuple(st, buf); err == nil {
			return Tuple{}, nil
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
// Bruijn index, a name off the stack a free variable the typechecker
// reports.
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

		return FreeVar{Name: name}, nil
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
	parseEmptyTuple := parsec.Try(parseRightParen)
	parseComma := parsec.Try(tokens.Exact(p.lexer, KindSymbol, ","))

	return func(st *state, buf parsec.Buffer[rune, Position]) (Pat, Error) {
		if _, err := parseLeftParen(st, buf); err != nil {
			return nil, err
		}

		// the unit pattern () is the empty tuple
		if _, err := parseEmptyTuple(st, buf); err == nil {
			return TuplePat{}, nil
		}

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
