package bytes

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/stretchr/testify/require"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	mul := func(x, y string) string {
		return fmt.Sprintf("(%s * %s)", x, y)
	}

	plus := func(x, y string) string {
		return fmt.Sprintf("(%s + %s)", x, y)
	}

	parseOp := func(buf parsec.Buffer[byte, int]) (parsec.BinaryOp[string], parsec.Error[int]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, parsec.NewParseError(pos, err.Error())
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

	parseOp := func(buf parsec.Buffer[byte, int]) (parsec.BinaryOp[string], parsec.Error[int]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, parsec.NewParseError(pos, err.Error())
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
					err:    parsec.NewParseError(0, "end of file"),
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

	parseOp := func(buf parsec.Buffer[byte, int]) (parsec.BinaryOp[string], parsec.Error[int]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, parsec.NewParseError(pos, err.Error())
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

	parseOp := func(buf parsec.Buffer[byte, int]) (parsec.BinaryOp[string], parsec.Error[int]) {
		pos := buf.Position()

		symbol, err := buf.Read(true)
		if err != nil {
			return nil, parsec.NewParseError(pos, err.Error())
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
					err:    parsec.NewParseError(0, "end of file"),
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
					output: []byte{},
				},
				{
					input:  []byte("a,b,c"),
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
					err:    parsec.NewParseError(0, "expected at least one item separated by ','"),
				},
				{
					input:  []byte("a,b,c"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte(","),
					output: nil,
					err:    parsec.NewParseError(0, "expected at least one item separated by ','"),
				},
				{
					input:  []byte(",a,b,c"),
					output: nil,
					err:    parsec.NewParseError(0, "expected at least one item separated by ','"),
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
					output: []byte{},
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
					output: []byte{},
				},
				{
					input:  []byte(",a,b,c"),
					output: []byte{},
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("a"),
					output: []byte{},
				},
				{
					input:  []byte("a,"),
					output: []byte{'a'},
				},
				{
					input:  []byte(",a"),
					output: []byte{},
				},
				{
					input:  []byte(",a,"),
					output: []byte{},
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
					input: []byte{},
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
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
					input: []byte(","),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input: []byte(",a,b,c"),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,b,c,"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input: []byte("a"),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input:  []byte("a,"),
					output: []byte{'a'},
				},
				{
					input: []byte(",a"),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input: []byte(",a,"),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
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
					output: []byte{},
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
					err:    parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
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
					err:    parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
				},
				{
					input: []byte(",a,b,c"),
					err:   parsec.NewParseError(0, "expected at least one item separated and ended by ','"),
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
					output: []byte{},
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
					input:  []byte(".1"),
					output: []byte{},
				},
				{
					input: []byte("a"),
					err:   parsec.NewParseError(0, "expected sequence of digits ended by dot"),
				},
				{
					input: []byte("123a."),
					err:   parsec.NewParseError(3, "expected sequence of digits ended by dot"),
				},
				{
					input:  []byte("123.a"),
					output: []byte{'1', '2', '3'},
				},
				{
					input: []byte("a123."),
					err:   parsec.NewParseError(0, "expected sequence of digits ended by dot"),
				},
				{
					input: []byte("12a3"),
					err:   parsec.NewParseError(2, "expected sequence of digits ended by dot"),
				},
			},
		},
	})
}

func TestSeq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input  string
		output []byte
		err    parsec.Error[int]
	}{
		{
			input:  "",
			output: nil,
		},
		{
			input:  "aa",
			output: []byte{'a', 'a'},
		},
		{
			input:  "b",
			output: nil,
			err:    parsec.NewParseError(0, "expected a"),
		},
		{
			input:  "aab",
			output: []byte{'a', 'a'},
			err:    parsec.NewParseError(2, "expected a"),
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			buf := Buffer([]byte(test.input))

			seq, err := Seq(Try(Eq("expected a", 'a')))(buf)
			require.NoError(t, err)

			var got []byte
			var last parsec.Error[int]

			for b, err := range seq {
				if err != nil {
					last = err

					break
				}

				got = append(got, b)
			}

			require.Equal(t, test.output, got)
			require.Equal(t, test.err, last)
		})
	}
}

func TestSeqEarlyStop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		take  int
		next  byte
	}{
		{
			input: "aab",
			take:  1,
			next:  'a',
		},
		{
			input: "aab",
			take:  2,
			next:  'b',
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s take %d", test.input, test.take), func(t *testing.T) {
			t.Parallel()

			buf := Buffer([]byte(test.input))

			seq, err := Seq(Try(Eq("expected a", 'a')))(buf)
			require.NoError(t, err)

			taken := 0

			for range seq {
				taken++

				if taken == test.take {
					break
				}
			}

			next, err := Try(Eq("expected next byte", test.next))(buf)
			require.NoError(t, err)
			require.Equal(t, test.next, next)
		})
	}
}
