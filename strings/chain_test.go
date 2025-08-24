package strings

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	mul := func(x, y string) string {
		return fmt.Sprintf("(%s * %s)", x, y)
	}

	plus := func(x, y string) string {
		return fmt.Sprintf("(%s + %s)", x, y)
	}

	parseOp := func(buf common.Buffer[rune, Position]) (common.BinaryOp[string], common.Error[Position]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, common.NewParseError(pos, err.Error())
		}

		if symbol == '+' {
			return plus, nil
		} else {
			return mul, nil
		}
	}

	parseItem := Cast(Any(), func(x rune) (string, error) {
		return fmt.Sprintf("%v", string(x)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainl("default", parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  "",
					output: "default",
				},
				{
					input:  "1+2*3?4",
					output: "(((1 + 2) * 3) * 4)",
				},
				{
					input:  "1+2*3",
					output: "((1 + 2) * 3)",
				},
				{
					input:  "1+2*3x",
					output: "((1 + 2) * 3)",
				},
				{
					input:  "1+2",
					output: "(1 + 2)",
				},
				{
					input:  "1*2",
					output: "(1 * 2)",
				},
				{
					input:  "1?2",
					output: "(1 * 2)",
				},
				{
					input:  "1+",
					output: "1",
				},
				{
					input:  "1",
					output: "1",
				},
			},
		},
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	mul := func(x, y string) string {
		return fmt.Sprintf("(%s * %s)", x, y)
	}

	plus := func(x, y string) string {
		return fmt.Sprintf("(%s + %s)", x, y)
	}

	parseOp := func(buf common.Buffer[rune, Position]) (common.BinaryOp[string], common.Error[Position]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, common.NewParseError(pos, err.Error())
		}

		if symbol == '+' {
			return plus, nil
		} else {
			return mul, nil
		}
	}

	parseItem := Cast(Any(), func(x rune) (string, error) {
		return fmt.Sprintf("%v", string(x)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainl1(parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  "",
					output: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"end of file",
					),
				},
				{
					input:  "1+2*3?4",
					output: "(((1 + 2) * 3) * 4)",
				},
				{
					input:  "1+2*3",
					output: "((1 + 2) * 3)",
				},
				{
					input:  "1+2*3x",
					output: "((1 + 2) * 3)",
				},
				{
					input:  "1+2",
					output: "(1 + 2)",
				},
				{
					input:  "1*2",
					output: "(1 * 2)",
				},
				{
					input:  "1?2",
					output: "(1 * 2)",
				},
				{
					input:  "1+",
					output: "1",
				},
				{
					input:  "1",
					output: "1",
				},
			},
		},
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	mul := func(x, y string) string {
		return fmt.Sprintf("(%s * %s)", x, y)
	}

	plus := func(x, y string) string {
		return fmt.Sprintf("(%s + %s)", x, y)
	}

	parseOp := func(buf common.Buffer[rune, Position]) (common.BinaryOp[string], common.Error[Position]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, common.NewParseError(pos, err.Error())
		}

		if symbol == '+' {
			return plus, nil
		} else {
			return mul, nil
		}
	}

	parseItem := Cast(Any(), func(x rune) (string, error) {
		return fmt.Sprintf("%v", string(x)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainr[string]("default", parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  "",
					output: "default",
				},
				{
					input:  "1+2*3?4",
					output: "(1 + (2 * (3 * 4)))",
				},
				{
					input:  "1+2*3",
					output: "(1 + (2 * 3))",
				},
				{
					input:  "1+2*3x",
					output: "(1 + (2 * 3))",
				},
				{
					input:  "1+2",
					output: "(1 + 2)",
				},
				{
					input:  "1*2",
					output: "(1 * 2)",
				},
				{
					input:  "1?2",
					output: "(1 * 2)",
				},
				{
					input:  "1+",
					output: "1",
				},
				{
					input:  "1",
					output: "1",
				},
			},
		},
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	mul := func(x, y string) string {
		return fmt.Sprintf("(%s * %s)", x, y)
	}

	plus := func(x, y string) string {
		return fmt.Sprintf("(%s + %s)", x, y)
	}

	parseOp := func(buf common.Buffer[rune, Position]) (common.BinaryOp[string], common.Error[Position]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, common.NewParseError(pos, err.Error())
		}

		if symbol == '+' {
			return plus, nil
		} else {
			return mul, nil
		}
	}

	parseItem := Cast(Any(), func(x rune) (string, error) {
		return fmt.Sprintf("%v", string(x)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainr1(parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  "",
					output: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"end of file",
					),
				},
				{
					input:  "1+2*3?4",
					output: "(1 + (2 * (3 * 4)))",
				},
				{
					input:  "1+2*3",
					output: "(1 + (2 * 3))",
				},
				{
					input:  "1+2*3x",
					output: "(1 + (2 * 3))",
				},
				{
					input:  "1+2",
					output: "(1 + 2)",
				},
				{
					input:  "1*2",
					output: "(1 * 2)",
				},
				{
					input:  "1?2",
					output: "(1 * 2)",
				},
				{
					input:  "1+",
					output: "1",
				},
				{
					input:  "1",
					output: "1",
				},
			},
		},
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: SepBy(
				0,
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: []rune{},
				},
				{
					input:  ",a,b,c",
					output: []rune{},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a",
					output: []rune{'a'},
				},
				{
					input:  "abc",
					output: []rune{'a'},
				},
			},
		},
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: SepBy1(
				0,
				"expected at least one item separated by ','",
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated by ','",
					),
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated by ','",
					),
				},
				{
					input:  ",a,b,c",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated by ','",
					),
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a",
					output: []rune{'a'},
				},
				{
					input:  "abc",
					output: []rune{'a'},
				},
			},
		},
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: EndBy(
				0,
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b'},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: []rune{},
				},
				{
					input:  ",a,b,c",
					output: []rune{},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a",
					output: []rune{},
				},
				{
					input:  "a,",
					output: []rune{'a'},
				},
				{
					input:  ",a",
					output: []rune{},
				},
				{
					input:  ",a,",
					output: []rune{},
				},
				{
					input:  "a,,",
					output: []rune{'a'},
				},
			},
		},
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: EndBy1(
				0,
				"expected at least one item separated and ended by ','",
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b'},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  ",a,b,c",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  "a,",
					output: []rune{'a'},
				},
				{
					input:  ",a",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  ",a,",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  "a,,",
					output: []rune{'a'},
				},
			},
		},
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: SepEndBy(
				0,
				NotEq("expected not eq ','", ','),
				Eq("expected eq ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,,d",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: []rune{},
				},
				{
					input:  ",a,b,c",
					output: []rune{},
				},
			},
		},
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: SepEndBy1(
				0,
				"expected at least one item separated and ended by ','",
				NotEq("expected not eq ','", ','),
				Eq("expected eq ','", ','),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  "a,b,c",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,,",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  "a,b,c,,d",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input:  ",",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
				{
					input:  ",a,b,c",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one item separated and ended by ','",
					),
				},
			},
		},
	})
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: ManyTill(
				0,
				"expected sequence of digits ended by dot",
				Range("expected digit", '0', '9'),
				Eq("expected dot", '.'),
			),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
				},
				{
					input:  "123.",
					output: []rune{'1', '2', '3'},
				},
				{
					input:  "123",
					output: []rune{'1', '2', '3'},
				},
				{
					input:  "123.45",
					output: []rune{'1', '2', '3'},
				},
				{
					input:  "1",
					output: []rune{'1'},
				},
				{
					input:  ".1",
					output: []rune{},
				},
				{
					input: "a",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected sequence of digits ended by dot",
					),
				},
				{
					input: "123a.",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 3,
							index:  3,
						},
						"expected sequence of digits ended by dot",
					),
				},
				{
					input:  "123.a",
					output: []rune{'1', '2', '3'},
				},
				{
					input: "a123.",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected sequence of digits ended by dot",
					),
				},
				{
					input: "12a3",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected sequence of digits ended by dot",
					),
				},
			},
		},
	})
}
