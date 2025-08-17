package bytes

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainl(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])
		c := 0

		comb := Chainl(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainl(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string { return "" }, common.NewParseError(buffer.Position(), "test error")
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainl1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])
		c := 0

		comb := Chainl1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainl1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "--", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string { return "++" }, common.NewParseError(buffer.Position(), "test error")
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainr(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])
		c := 0

		comb := Chainr(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainr(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string { return "" }, common.NewParseError(buffer.Position(), "test error")
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainr1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])
		c := 0

		comb := Chainr1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("any byte", true, common.Anything[byte])

		comb := Chainr1(
			func(buffer common.Buffer[byte, int]) (string, common.Error[int]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[byte, int]) (func(string, string) string, common.Error[int]) {
				return func(x, y string) string {
					return ""
				}, common.NewParseError(buffer.Position(), "test error")
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
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
