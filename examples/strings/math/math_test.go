package main

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	expr      string
	wantAST   Expr
	wantEval  int
	wantError bool
}

func TestMathParser(t *testing.T) {
	tests := []testCase{
		{
			expr: "1 + 2 * 3",
			wantAST: BinaryOp{
				Op:   '+',
				Left: Number{Value: 1},
				Right: BinaryOp{
					Op:    '*',
					Left:  Number{Value: 2},
					Right: Number{Value: 3},
				},
			},
			wantEval: 7,
		},
		{
			expr: "(1 + 2) * 3",
			wantAST: BinaryOp{
				Op: '*',
				Left: BinaryOp{
					Op:    '+',
					Left:  Number{Value: 1},
					Right: Number{Value: 2},
				},
				Right: Number{Value: 3},
			},
			wantEval: 9,
		},
		{
			expr: "4 + 5 * (6 - 7)",
			wantAST: BinaryOp{
				Op:   '+',
				Left: Number{Value: 4},
				Right: BinaryOp{
					Op:   '*',
					Left: Number{Value: 5},
					Right: BinaryOp{
						Op:    '-',
						Left:  Number{Value: 6},
						Right: Number{Value: 7},
					},
				},
			},
			wantEval: -1,
		},
		{
			expr: "10 / 2 + 3 * 2",
			wantAST: BinaryOp{
				Op: '+',
				Left: BinaryOp{
					Op:    '/',
					Left:  Number{Value: 10},
					Right: Number{Value: 2},
				},
				Right: BinaryOp{
					Op:    '*',
					Left:  Number{Value: 3},
					Right: Number{Value: 2},
				},
			},
			wantEval: 11,
		},
		{
			expr: "(1 + 2) * 3",
			wantAST: BinaryOp{
				Op: '*',
				Left: BinaryOp{
					Op:    '+',
					Left:  Number{Value: 1},
					Right: Number{Value: 2},
				},
				Right: Number{Value: 3},
			},
			wantEval: 9,
		},
		{
			expr:     "  42",
			wantAST:  Number{Value: 42},
			wantEval: 42,
		},
		{
			expr:     "(((7)))",
			wantAST:  Number{Value: 7},
			wantEval: 7,
		},
		{
			expr:      "1 +",
			wantError: true,
		},
		{
			expr:      "2 * (3 + )",
			wantError: true,
		},
		{
			expr:      "abc",
			wantError: true,
		},
		{
			expr:      "1 + 2 2",
			wantError: true,
		},
		{
			expr:      "",
			wantError: true,
		},
	}

	for i, x := range tests {
		tc := x

		name := fmt.Sprintf("%d - %s", i, tc.expr)

		t.Run(name, func(t *testing.T) {
			parse := Parser(t)
			buf := strings.Buffer([]rune(tc.expr))
			ast, err := parse(buf)

			if tc.wantError {
				require.Error(t, err, "expected error for input: %q", tc.expr)
				return
			}

			require.NoError(t, err, "unexpected error for input: %q", tc.expr)
			require.True(t, buf.IsEOF(), "expected all input to be consumed for: %q", tc.expr)
			require.NotNil(t, ast, "expected AST for: %q", tc.expr)
			require.Equal(t, tc.wantAST, ast, "AST mismatch for: %q", tc.expr)
			require.Equal(t, tc.wantEval, ast.Eval(), "Eval mismatch for: %q", tc.expr)
		})
	}
}
