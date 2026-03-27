package main

import (
	"fmt"
)

// Expr — AST для математических выражений
type Expr interface {
	Eval() int
	String() string
}

type Number struct {
	Value int
}

func (n Number) Eval() int      { return n.Value }
func (n Number) String() string { return fmt.Sprintf("%d", n.Value) }

type BinaryOp struct {
	Left  Expr
	Op    rune
	Right Expr
}

func (b BinaryOp) Eval() int {
	switch b.Op {
	case '+':
		return b.Left.Eval() + b.Right.Eval()
	case '-':
		return b.Left.Eval() - b.Right.Eval()
	case '*':
		return b.Left.Eval() * b.Right.Eval()
	case '/':
		return b.Left.Eval() / b.Right.Eval()
	default:
		panic("unknown op")
	}
}
func (b BinaryOp) String() string {
	return fmt.Sprintf("(%s %c %s)", b.Left.String(), b.Op, b.Right.String())
}
