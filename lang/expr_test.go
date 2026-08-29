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
