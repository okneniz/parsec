package main

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/okneniz/parsec/common"
	"github.com/okneniz/parsec/strings"
)

func parseNumber() common.Combinator[rune, strings.Position, Expr] {
	parse := strings.Padded(
		strings.Try(strings.Space("whitespaces")),
		strings.Unsigned[int](),
	)

	var null Expr

	return func(
		buf common.Buffer[rune, strings.Position],
	) (Expr, common.Error[strings.Position]) {
		val, err := parse(buf)
		if err != nil {
			return null, err
		}

		return Number{Value: val}, nil
	}
}

func Parser(t testing.TB) common.Combinator[rune, strings.Position, Expr] {
	parseNum := strings.Try(parseNumber())

	parseLeftParens := strings.Padded(
		strings.Try(strings.Space("whitespaces")),
		strings.Eq("expectef left parens", '('),
	)

	parseRightParens := strings.Padded(
		strings.Try(strings.Space("whitespaces")),
		strings.Eq("expectef right parens", ')'),
	)

	prioritiesBinary := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	maxPriority := 2

	parseOperation := strings.Try(strings.OneOf(
		"expected one of operation",
		slices.Collect(maps.Keys(prioritiesBinary))...,
	))

	parseOperation = strings.Padded(
		strings.Try(strings.Space("whitespaces")),
		parseOperation,
	)

	var parseOperand common.Combinator[rune, strings.Position, Expr]
	var parseBinaryOp func(int) common.Combinator[rune, strings.Position, Expr]
	var parseExpr common.Combinator[rune, strings.Position, Expr]

	parseOperand = strings.Try(func(
		buf common.Buffer[rune, strings.Position],
	) (Expr, common.Error[strings.Position]) {
		val, err := parseNum(buf)
		if err == nil {
			return val, nil
		}

		_, err = parseLeftParens(buf)
		if err != nil {
			return nil, err
		}

		expr, err := parseExpr(buf)
		if err != nil {
			return nil, err
		}

		_, err = parseRightParens(buf)
		if err != nil {
			return nil, err
		}

		return expr, nil
	})

	priorityParsers := make(map[int]common.Combinator[rune, strings.Position, Expr])

	parseBinaryOp = func(prevPriority int) common.Combinator[rune, strings.Position, Expr] {
		if parse, exists := priorityParsers[prevPriority]; exists {
			return parse
		}

		if prevPriority > maxPriority {
			return func(
				buf common.Buffer[rune, strings.Position],
			) (Expr, common.Error[strings.Position]) {
				return parseOperand(buf)
			}
		}

		parse := func(
			buf common.Buffer[rune, strings.Position],
		) (Expr, common.Error[strings.Position]) {
			left, err := parseOperand(buf)
			if err != nil {
				return nil, err
			}

			operation, err := parseOperation(buf)
			if err != nil {
				return left, nil
			}

			priority, exists := prioritiesBinary[operation]
			if !exists {
				return nil, common.NewParseError(
					buf.Position(),
					fmt.Sprintf("unknown operation: %v", string(operation)),
				)
			}

			if prevPriority > priority {
				return nil, common.NewParseError(
					buf.Position(),
					fmt.Sprintf("nothing matched: %v", string(operation)),
				)
			}

			nextExpression, exists := priorityParsers[priority+1]
			if !exists {
				nextExpression = parseOperand
			}

			right, err := nextExpression(buf) // higher priority or value
			if err != nil {
				return nil, err
			}

			return BinaryOp{
				Left:  left,
				Op:    operation,
				Right: right,
			}, nil
		}

		return parse
	}

	for i := 0; i <= maxPriority; i++ {
		priorityParsers[i] = parseBinaryOp(i)
	}

	parseExpr = func(
		buf common.Buffer[rune, strings.Position],
	) (Expr, common.Error[strings.Position]) {
		left, err := priorityParsers[0](buf)
		if err != nil {
			return nil, err
		}

		operation, err := parseOperation(buf)
		if err != nil {
			return left, nil
		}

		priority, exists := prioritiesBinary[operation]
		if !exists {
			return nil, common.NewParseError(
				buf.Position(),
				fmt.Sprintf("unknown operation: %v", string(operation)),
			)
		}

		if 0 > priority {
			return nil, common.NewParseError(
				buf.Position(),
				fmt.Sprintf("nothing matched: %v", string(operation)),
			)
		}

		nextExpression, exists := priorityParsers[priority+1]
		if !exists {
			nextExpression = parseNum
		}

		right, err := nextExpression(buf) // higher priority or value
		if err != nil {
			return nil, err
		}

		return BinaryOp{
			Left:  left,
			Op:    operation,
			Right: right,
		}, nil
	}

	eof := strings.Padded(
		strings.Try(strings.Space("whitespaces")),
		strings.EOF(),
	)

	return func(
		buf common.Buffer[rune, strings.Position],
	) (Expr, common.Error[strings.Position]) {
		expression, err := parseExpr(buf)
		if err != nil {
			return nil, err
		}

		if ok, _ := eof(buf); !ok {
			return nil, common.NewParseError(
				buf.Position(),
				"expected end of file",
			)
		}

		return expression, nil
	}
}
