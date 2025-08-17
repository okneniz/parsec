package bytes

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

	parseOp := func(buf common.Buffer[byte, int]) (common.BinaryOp[string], common.Error[int]) {
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

	parseItem := Cast(Any(), func(b byte) (string, error) {
		return fmt.Sprintf("%d", int(b)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainl("default", parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  []byte{},
					output: "default",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, '?', 4},
					output: "(((1 + 2) * 3) * 4)",
				},
				{
					input:  []byte{1, '+', 2, '*', 3},
					output: "((1 + 2) * 3)",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, 100},
					output: "((1 + 2) * 3)",
				},
				{
					input:  []byte{1, '+', 2},
					output: "(1 + 2)",
				},
				{
					input:  []byte{1, '*', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '?', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '+'},
					output: "1",
				},
				{
					input:  []byte{1},
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

	parseOp := func(buf common.Buffer[byte, int]) (common.BinaryOp[string], common.Error[int]) {
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

	parseItem := Cast(Any(), func(b byte) (string, error) {
		return fmt.Sprintf("%d", int(b)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainl1(parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  []byte{},
					output: "",
					err:    common.NewParseError(0, "end of file"),
				},
				{
					input:  []byte{1, '+', 2, '*', 3, '?', 4},
					output: "(((1 + 2) * 3) * 4)",
				},
				{
					input:  []byte{1, '+', 2, '*', 3},
					output: "((1 + 2) * 3)",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, 100},
					output: "((1 + 2) * 3)",
				},
				{
					input:  []byte{1, '+', 2},
					output: "(1 + 2)",
				},
				{
					input:  []byte{1, '*', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '?', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '+'},
					output: "1",
				},
				{
					input:  []byte{1},
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

	parseOp := func(buf common.Buffer[byte, int]) (common.BinaryOp[string], common.Error[int]) {
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

	parseItem := Cast(Any(), func(b byte) (string, error) {
		return fmt.Sprintf("%d", int(b)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainr("default", parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  []byte{},
					output: "default",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, '?', 4},
					output: "(1 + (2 * (3 * 4)))",
				},
				{
					input:  []byte{1, '+', 2, '*', 3},
					output: "(1 + (2 * 3))",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, 100},
					output: "(1 + (2 * 3))",
				},
				{
					input:  []byte{1, '+', 2},
					output: "(1 + 2)",
				},
				{
					input:  []byte{1, '*', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '?', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '+'},
					output: "1",
				},
				{
					input:  []byte{1},
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

	parseOp := func(buf common.Buffer[byte, int]) (common.BinaryOp[string], common.Error[int]) {
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

	parseItem := Cast(Any(), func(b byte) (string, error) {
		return fmt.Sprintf("%d", int(b)), nil
	})

	runTests(t, []test[string]{
		{
			comb: Chainr1(parseItem, parseOp),
			cases: []testCase[string]{
				{
					input:  []byte{},
					output: "",
					err:    common.NewParseError(0, "end of file"),
				},
				{
					input:  []byte{1, '+', 2, '*', 3, '?', 4},
					output: "(1 + (2 * (3 * 4)))",
				},
				{
					input:  []byte{1, '+', 2, '*', 3},
					output: "(1 + (2 * 3))",
				},
				{
					input:  []byte{1, '+', 2, '*', 3, 100},
					output: "(1 + (2 * 3))",
				},
				{
					input:  []byte{1, '+', 2},
					output: "(1 + 2)",
				},
				{
					input:  []byte{1, '*', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '?', 2},
					output: "(1 * 2)",
				},
				{
					input:  []byte{1, '+'},
					output: "1",
				},
				{
					input:  []byte{1},
					output: "1",
				},
			},
		},
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: SepBy(
				0,
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
				},
				{
					input:  []byte(",a,b,c"),
					output: nil,
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a"),
					output: []byte{'a'},
				},
				{
					input:  []byte("abc"),
					output: []byte{'a'},
				},
			},
		},
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: SepBy1(
				0,
				"expected at least one item separated by ','",
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated by ','"),
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated by ','"),
				},
				{
					input:  []byte(",a,b,c"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated by ','"),
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a"),
					output: []byte{'a'},
				},
				{
					input:  []byte("abc"),
					output: []byte{'a'},
				},
			},
		},
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: EndBy(
				0,
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b'},
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
				},
				{
					input:  []byte(",a,b,c"),
					output: nil,
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a"),
					output: nil,
				},
				{
					input:  []byte("a,"),
					output: []byte{'a'},
				},
				{
					input:  []byte(",a"),
					output: nil,
				},
				{
					input:  []byte(",a,"),
					output: nil,
				},
				{
					input:  []byte("a,,"),
					output: []byte{'a'},
				},
			},
		},
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: EndBy1(
				0,
				"expected at least one item separated and ended by ','",
				NotEq("expected not ','", ','),
				Eq("expected ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b'},
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte(",a,b,c"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,"),
					output: []byte{'a'},
				},
				{
					input:  []byte(",a"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte(",a,"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,,"),
					output: []byte{'a'},
				},
			},
		},
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: SepEndBy(
				0,
				NotEq("expected not eq ','", ','),
				Eq("expected eq ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,,d"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: []byte{},
				},
				{
					input:  []byte(",a,b,c"),
					output: []byte{},
				},
			},
		},
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: SepEndBy1(
				0,
				"expected at least one item separated and ended by ','",
				NotEq("expected not eq ','", ','),
				Eq("expected eq ','", ','),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a,b,c,,d"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
					err:    common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input: []byte(",a,b,c"),
					err:   common.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
			},
		},
	})
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: ManyTill(
				0,
				"expected sequence of digits ended by dot",
				Range("expected digit", '0', '9'),
				Eq("expected dot", '.'),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
				},
				{
					input:  []byte("123."),
					output: []byte{'1', '2', '3'},
				},
				{
					input:  []byte("123"),
					output: []byte{'1', '2', '3'},
				},
				{
					input:  []byte("123.45"),
					output: []byte{'1', '2', '3'},
				},
				{
					input:  []byte("1"),
					output: []byte{'1'},
				},
				{
					input: []byte(".1"),
					err:   nil,
				},
				{
					input: []byte("a"),
					err:   common.NewParseError(0, "expected sequence of digits ended by dot"),
				},
				{
					input: []byte("123a."),
					err:   common.NewParseError(3, "expected sequence of digits ended by dot"),
				},
				{
					input:  []byte("123.a"),
					output: []byte{'1', '2', '3'},
				},
				{
					input: []byte("a123."),
					err:   common.NewParseError(0, "expected sequence of digits ended by dot"),
				},
				{
					input: []byte("12a3"),
					err:   common.NewParseError(2, "expected sequence of digits ended by dot"),
				},
			},
		},
	})
}
