package lang

import (
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
	"github.com/stretchr/testify/require"
)

func TestExpr(t *testing.T) {
	p := NewParser().
		Infix(1, "+", "-").
		Infix(2, "*", "/").
		InfixRight(3, "^").
		Prefix("-").
		Postfix("!")

	def := testDefinition()
	def.Operators = append(def.Operators, "^", "!")

	tests := []struct {
		src  string
		want string
	}{
		{"1", "1"},
		{"x", "x"},
		{"1 + 2", "(+ 1 2)"},
		{"1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"f x y", "(app (app f x) y)"},
		{"1 2", "(app 1 2)"},
		{"f -x", "(- f x)"},
		{"1 * 2 + 3", "(+ (* 1 2) 3)"},
		{"a - b - c", "(- (- a b) c)"},
		{"a ^ b ^ c", "(^ a (^ b c))"},
		{"(1 + 2) * 3", "(* (+ 1 2) 3)"},
		{"-x + y", "(+ (- x) y)"},
		{"-(x + y)", "(- (+ x y))"},
		{"2 - -3", "(- 2 (- 3))"},
		{"1 + 2 // comment\n * 3", "(+ 1 (* 2 3))"},
		{"a!", "(! a)"},
		{"a!!", "(! (! a))"},
		{"a! + b!", "(+ (! a) (! b))"},
		{"(a + b)!", "(! (+ a b))"},
		{"-a!", "(- (! a))"},
		{"1 (* a (* b *) c *) + 2", "(+ 1 2)"},
	}

	parse := New(def, p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}
}

func TestExprErrors(t *testing.T) {
	p := NewParser().Infix(1, "+", "-")

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"dangling operator", "1 +", "expected operand after operator"},
		{"unclosed paren", "(1 + 2", "expected operand"},
		{"trailing garbage", "1 )", "expected end of input"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(testDefinition(), p)(strings.Buffer([]rune(test.src)))
			require.ErrorContains(t, err, test.want)
		})
	}
}

// letNode stands in for a real binding node in this test.
type letNode struct {
	Name  string
	Value Expr
	Body  Expr
}

func (n letNode) String() string {
	return "(let " + n.Name + " " + n.Value.String() + " " + n.Body.String() + ")"
}

func TestExtraOperands(t *testing.T) {
	p := NewParser().Infix(1, "+", "-").Infix(2, "*", "/")

	// let x = value in body, with the value and the body parsed by
	// the same parser — the closure-over-parser pattern.
	let := func(lex tokens.Lexer[Kind, Lexeme]) parsec.Combinator[rune, strings.Position, Expr] {
		letKw := tokens.Exact(lex, KindKeyword, "let")
		inKw := tokens.Exact(lex, KindKeyword, "in")
		eq := tokens.Exact(lex, KindSymbol, "=")
		ident := tokens.OfKind(lex, KindIdent)

		return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
			if _, err := letKw(buf); err != nil {
				return nil, err
			}

			name, err := ident(buf)
			if err != nil {
				return nil, err
			}

			if _, err := eq(buf); err != nil {
				return nil, err
			}

			value, verr := p.Expr(lex)(buf)
			if verr != nil {
				return nil, verr
			}

			if _, err := inKw(buf); err != nil {
				return nil, err
			}

			body, berr := p.Expr(lex)(buf)
			if berr != nil {
				return nil, berr
			}

			return letNode{Name: string(name.Lexeme), Value: value, Body: body}, nil
		}
	}

	e, err := New(testDefinition(), p, let)(strings.Buffer([]rune("let x = 1 + 2 in x * 3")))
	require.NoError(t, err)
	require.Equal(t, "(let x (+ 1 2) (* x 3))", e.String())

	// extras are tried first, with backtracking: a plain operand
	// still parses when the extra declines.
	e, err = New(testDefinition(), p, let)(strings.Buffer([]rune("1 + 2")))
	require.NoError(t, err)
	require.Equal(t, "(+ 1 2)", e.String())
}

// sizeNode stands in for a keyword-led operand form of this test.
type sizeNode struct {
	Sub Expr
}

func (n sizeNode) String() string {
	return "(sizeof " + n.Sub.String() + ")"
}

// castNode stands in for a symbol-led operand form of this test.
type castNode struct {
	Type string
	Sub  Expr
}

func (n castNode) String() string {
	return "(cast " + n.Type + " " + n.Sub.String() + ")"
}

func TestPrefixForm(t *testing.T) {
	sizeof := func(l Ladder) parsec.Combinator[rune, strings.Position, Expr] {
		kw := tokens.Exact(l.Lex, KindKeyword, "sizeof")

		return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
			if _, err := kw(buf); err != nil {
				return nil, err
			}

			sub, serr := l.Atom(buf)
			if serr != nil {
				return nil, serr
			}

			return sizeNode{Sub: sub}, nil
		}
	}

	cast := func(l Ladder) parsec.Combinator[rune, strings.Position, Expr] {
		open := tokens.Exact(l.Lex, KindSymbol, "(")
		closeParen := tokens.Exact(l.Lex, KindSymbol, ")")
		ident := tokens.OfKind(l.Lex, KindIdent)

		return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
			if _, err := open(buf); err != nil {
				return nil, err
			}

			ty, terr := ident(buf)
			if terr != nil {
				return nil, terr
			}

			if _, err := closeParen(buf); err != nil {
				return nil, err
			}

			sub, serr := l.Atom(buf)
			if serr != nil {
				return nil, serr
			}

			return castNode{Type: string(ty.Lexeme), Sub: sub}, nil
		}
	}

	p := NewParser().
		Infix(1, "+", "-").
		Infix(2, "*", "/").
		PrefixForm(sizeof, cast)

	def := testDefinition()
	def.Keywords = append(def.Keywords, "sizeof")

	tests := []struct {
		src  string
		want string
	}{
		{"sizeof x", "(sizeof x)"},
		{"sizeof x * 2", "(* (sizeof x) 2)"},
		{"1 + sizeof y", "(+ 1 (sizeof y))"},
		{"sizeof sizeof x", "(sizeof (sizeof x))"},
		{"( int ) x * 2", "(* (cast int x) 2)"},
		{"(x + 2) * 3", "(* (+ x 2) 3)"},
		{"( int ) sizeof x", "(cast int (sizeof x))"},
	}

	parse := New(def, p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}
}

// ifNode stands in for an expression form of this test.
type ifNode struct {
	Cond, Then, Else Expr
}

func (n ifNode) String() string {
	return "(if " + n.Cond.String() + " " + n.Then.String() + " " + n.Else.String() + ")"
}

func TestExprForm(t *testing.T) {
	ifForm := func(l Ladder) parsec.Combinator[rune, strings.Position, Expr] {
		ifKw := tokens.Exact(l.Lex, KindKeyword, "if")
		thenKw := tokens.Exact(l.Lex, KindKeyword, "then")
		elseKw := tokens.Exact(l.Lex, KindKeyword, "else")

		return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
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

			return ifNode{Cond: cond, Then: then, Else: els}, nil
		}
	}

	p := NewParser().
		Infix(1, "+", "-").
		Infix(2, "*", "/").
		ExprForm(ifForm)

	def := testDefinition()
	def.Keywords = append(def.Keywords, "if", "then", "else")

	tests := []struct {
		src  string
		want string
	}{
		{"if x then 1 else 2", "(if x 1 2)"},
		{"if x then 1 + 2 else 3 * 4", "(if x (+ 1 2) (* 3 4))"},
		{"if a then if b then 1 else 2 else 3", "(if a (if b 1 2) 3)"},
		{"(if x then 1 else 2) * 3", "(* (if x 1 2) 3)"},
		{"f (if x then y else z)", "(app f (if x y z))"},
	}

	parse := New(def, p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}

	// a form is not an operand: an if does not continue an infix
	// chain, and the whole expression fails
	_, err := parse(strings.Buffer([]rune("1 + if x then 1 else 2")))
	require.ErrorContains(t, err, "expected expression")
}

func TestKeywordInfix(t *testing.T) {
	p := NewParser().
		Infix(2, "*", "/").
		Infix(1, "+", "-").
		InfixRight(0, "andalso").
		InfixRight(-1, "orelse")

	def := testDefinition()
	def.Keywords = append(def.Keywords, "andalso", "orelse")

	tests := []struct {
		src  string
		want string
	}{
		{"a andalso b", "(andalso a b)"},
		{"a andalso b andalso c", "(andalso a (andalso b c))"},
		{"a orelse b andalso c", "(orelse a (andalso b c))"},
		{"a andalso b + c", "(andalso a (+ b c))"},
		{"a + b orelse c * d", "(orelse (+ a b) (* c d))"},
	}

	parse := New(def, p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}
}

// ternaryDefinition is the test definition with the two separators
// of the conditional among its symbols.
func ternaryDefinition() Definition {
	def := testDefinition()
	def.Punctuation = def.Punctuation + "?:"

	return def
}

func TestTernary(t *testing.T) {
	// the operator ladder of c: a comma of the lowest precedence,
	// the assignment above it, the conditional between the
	// assignment and the additive operators
	p := NewParser().
		Infix(0, ",").
		InfixRight(1, "=").
		Ternary(2, "?", ":").
		Infix(3, "+", "-").
		Infix(4, "*", "/")

	tests := []struct {
		src  string
		want string
	}{
		{"a ? b : c", "(?: a b c)"},
		{"a ? b : c ? d : e", "(?: a b (?: c d e))"},
		{"x = a ? b : c", "(= x (?: a b c))"},
		{"x = a + 2 ? b : c", "(= x (?: (+ a 2) b c))"},
		{"x = y = a + 1 ? b : c", "(= x (= y (?: (+ a 1) b c)))"},
		{"a + 1 = b ? c : d", "(= (+ a 1) (?: b c d))"},
		{"x = y = 1", "(= x (= y 1))"},
		{"1 + a ? b : c", "(?: (+ 1 a) b c)"},
		{"(a ? b : c) + 1", "(+ (?: a b c) 1)"},
		{"a ? b = 1 : c", "(?: a (= b 1) c)"},
		{"a ? b, c : d", "(?: a (, b c) d)"},
		{"a ? b : c, d", "(, (?: a b c) d)"},
		{"a, b, c", "(, (, a b) c)"},
		{"(a, b)", "(, a b)"},
	}

	parse := New(ternaryDefinition(), p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}
}

func TestTernaryErrors(t *testing.T) {
	p := NewParser().
		InfixRight(1, "=").
		Ternary(2, "?", ":").
		Infix(3, "+", "-")

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"missing colon", "a ? b", "to close the conditional operator"},
		{"missing middle", "a ? : c", "expected operand"},
		{"missing condition", "? a : b", "expected operand"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(ternaryDefinition(), p)(strings.Buffer([]rune(test.src)))
			require.ErrorContains(t, err, test.want)
		})
	}
}

// callNode stands in for a call form of this test.
type callNode struct {
	Fn   string
	Args []Expr
}

func (n callNode) String() string {
	out := "(" + n.Fn

	for _, arg := range n.Args {
		out += " " + arg.String()
	}

	return out + ")"
}

func TestLadderFrom(t *testing.T) {
	// a call reads its arguments above the comma, the assignment
	// still inside: the c argument list
	call := func(l Ladder) parsec.Combinator[rune, strings.Position, Expr] {
		open := tokens.Exact(l.Lex, KindSymbol, "(")
		closeParen := tokens.Exact(l.Lex, KindSymbol, ")")
		comma := parsec.Try(tokens.Exact(l.Lex, KindSymbol, ","))
		ident := tokens.OfKind(l.Lex, KindIdent)
		arg := l.From(1)

		return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
			fn, err := ident(buf)
			if err != nil {
				return nil, err
			}

			if _, err := open(buf); err != nil {
				return nil, err
			}

			node := callNode{Fn: string(fn.Lexeme)}

			for {
				e, aerr := arg(buf)
				if aerr != nil {
					return nil, aerr
				}

				node.Args = append(node.Args, e)

				if _, cerr := comma(buf); cerr != nil {
					break
				}
			}

			if _, err := closeParen(buf); err != nil {
				return nil, err
			}

			return node, nil
		}
	}

	p := NewParser().
		Infix(0, ",").
		InfixRight(1, "=").
		Ternary(2, "?", ":").
		Infix(3, "+", "-").
		PrefixForm(call)

	tests := []struct {
		src  string
		want string
	}{
		{"f(a, b)", "(f a b)"},
		{"f(a, b, c)", "(f a b c)"},
		{"g((a, b))", "(g (, a b))"},
		{"h(x = 1, y = 2)", "(h (= x 1) (= y 2))"},
		{"f(a ? b : c, d)", "(f (?: a b c) d)"},
		{"f(a, b), f(c, d)", "(, (f a b) (f c d))"},
	}

	parse := New(ternaryDefinition(), p)

	for _, test := range tests {
		expression := strings.Buffer([]rune(test.src))

		t.Run(test.src, func(t *testing.T) {
			e, err := parse(expression)
			require.NoError(t, err)
			require.Equal(t, test.want, e.String())
		})
	}
}
