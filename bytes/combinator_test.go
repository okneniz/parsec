package bytes

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

type (
	testCase[T any] struct {
		input  []byte
		output T
		err    error
	}

	test[T any] struct {
		comb  common.Combinator[byte, int, T]
		cases []testCase[T]
	}
)

func runTests[T comparable](t *testing.T, tests []test[T]) {
	t.Helper()

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("test %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, x := range test.cases {
				testCase := x
				name := fmt.Sprintf("case %d", i)

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					result, err := Parse(testCase.input, test.comb)

					if testCase.err != nil {
						AssertError(t, err)
						AssertEq(t, err.Error(), testCase.err.Error())
					} else {
						Check(t, err)
					}

					AssertEq(t, result, testCase.output)
				})
			}
		})
	}
}

func runTestsSlice[T comparable](t *testing.T, tests []test[[]T]) {
	t.Helper()

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("test %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, x := range test.cases {
				testCase := x
				name := fmt.Sprintf("case %d", i)

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					result, err := Parse(testCase.input, test.comb)

					if testCase.err != nil {
						AssertError(t, err)
						AssertEq(t, err.Error(), testCase.err.Error())
					} else {
						Check(t, err)
					}

					AssertSlice(t, result, testCase.output)
				})
			}
		})
	}
}

func TestSatisfy(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Satisfy("expected not 'c'", true, func(x byte) bool { return x != 'c' }),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected not 'c'"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("b"),
					output: 'b',
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected not 'c'"),
				},
			},
		},
		{
			comb: Satisfy("error explanation", true, func(x byte) bool { return false }),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
				{
					input:  []byte("abc"),
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
			},
		},
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Any(),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, common.ErrEndOfFile.Error()),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
			},
		},
	})

	t.Run("must parse first element from none empty input", func(t *testing.T) {
		t.Parallel()

		seed := time.Now().UnixNano()
		t.Logf("seed: %v", seed)
		rnd := rand.New(rand.NewPCG(0, uint64(seed)))
		comb := Any()

		ohsnap.Check(t, 1000, ohsnap.NewBuilder(rnd).Byte(), func(b byte) bool {
			result, err := Parse([]byte{b}, comb)
			if err != nil {
				t.Logf("input: %v", []byte{b})
				t.Logf("output: %v", result)
				t.Error(err)
				return false
			}

			return result == b
		})
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Try(
				Satisfy(
					"error explanation",
					true,
					func(x byte) bool { return x <= byte('b') },
				),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("b"),
					output: 'b',
				},
				{
					input: []byte("c"),
					err:   common.NewParseError(0, "error explanation"),
				},
			},
		},
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	notBrackets := Satisfy("test", true, func(x byte) bool {
		return !(x == byte(')') || x == byte('('))
	})

	comb := Between(
		Eq("expected '('", '('),
		Some(0, "expected not ( or ) symbols", Try(notBrackets)),
		Eq("expected ')'", ')'),
	)

	runTestsSlice(t, []test[[]byte]{
		{
			comb: comb,
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected '('"),
				},
				{
					input:  []byte("(abc)"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input: []byte("abc"),
					err:   common.NewParseError(0, "expected '('"),
				},
				{
					input: []byte("(abc"),
					err:   common.NewParseError(4, "expected ')'"),
				},
				{
					input: []byte(" (abc) "),
					err:   common.NewParseError(0, "expected '('"),
				},
				{
					input: []byte("((abc))"),
					err:   common.NewParseError(1, "expected not ( or ) symbols"),
				},
				{
					input: []byte("()"),
					err:   common.NewParseError(1, "expected not ( or ) symbols"),
				},
			},
		},
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("skip not optional", func(t *testing.T) {
		runTests(t, []test[byte]{
			{
				comb: Skip(
					Eq("expected 'a'", 'a'),
					Eq("expected 'b'", 'b'),
				),
				cases: []testCase[byte]{
					{
						input:  []byte{},
						output: 0,
						err:    common.NewParseError(0, "expected 'a'"),
					},
					{
						input:  []byte("abc"),
						output: byte('b'),
					},
					{
						input:  []byte("b"),
						output: 0,
						err:    common.NewParseError(0, "expected 'a'"),
					},
					{
						input: []byte("bbb"),
						err:   common.NewParseError(0, "expected 'a'"),
					},
					{
						input:  []byte("ac"),
						output: 0,
						err:    common.NewParseError(1, "expected 'b'"),
					},
				},
			},
		})
	})

	t.Run("skip optional", func(t *testing.T) {
		runTests(t, []test[byte]{
			{
				comb: Skip(
					Optional(Try(Eq("expected 'a'", 'a')), 'x'),
					Eq("expected 'b'", 'b'),
				),
				cases: []testCase[byte]{
					{
						input:  []byte{},
						output: 0,
						err:    common.NewParseError(0, "expected 'b'"),
					},
					{
						input:  []byte("abc"),
						output: byte('b'),
					},
					{
						input:  []byte("b"),
						output: 'b',
					},
					{
						input:  []byte("bbb"),
						output: 'b',
					},
					{
						input:  []byte("ac"),
						output: 0,
						err:    common.NewParseError(1, "expected 'b'"),
					},
				},
			},
		})
	})
}

func TestSkipMany(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: common.SkipMany(
				NoneOf("expected not a, b or c", 'a', 'b', 'c'),
				SequenceOf("expected abc", 'a', 'b', 'c'),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected abc"),
				},
				{
					input:  []byte("abc"),
					output: []byte("abc"),
				},
				{
					input:  []byte("ab"),
					output: nil,
					err:    common.NewParseError(0, "expected abc"),
				},
				{
					input:  []byte("xab"),
					output: nil,
					err:    common.NewParseError(1, "expected abc"),
				},
				{
					input:  []byte("123abc"),
					output: []byte("abc"),
				},
				{
					input:  []byte("bcabc"),
					output: nil,
					err:    common.NewParseError(0, "expected abc"),
				},
				{
					input:  []byte("abcabc"),
					output: []byte("abc"),
				},
				{
					input:  []byte("123abc123"),
					output: []byte("abc"),
				},
				{
					input:  []byte("123abcabc"),
					output: []byte("abc"),
				},
				{
					input:  []byte("123"),
					output: nil,
					err:    common.NewParseError(3, "expected abc"),
				},
			},
		},
	})
}

func TestSkipAfter(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: SkipAfter(
				Eq("expected 'b'", 'b'),
				Eq("expected 'a'", 'a'),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected 'a'"),
				},
				{
					input:  []byte("abc"),
					output: 'a',
				},
				{
					input:  []byte("ab"),
					output: 'a',
				},
				{
					input: []byte("a"),
					err:   common.NewParseError(1, "expected 'b'"),
				},
				{
					input: []byte("ac"),
					err:   common.NewParseError(1, "expected 'b'"),
				},
				{
					input: []byte("b"),
					err:   common.NewParseError(0, "expected 'a'"),
				},
				{
					input: []byte("bc"),
					err:   common.NewParseError(0, "expected 'a'"),
				},
				{
					input: []byte("bb"),
					err:   common.NewParseError(0, "expected 'a'"),
				},
			},
		},
	})
}

func TestPadded(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Padded(
				Eq("expected dot", '.'),
				Range("expected digit", '0', '9'),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected digit"),
				},
				{
					input:  []byte("1"),
					output: '1',
				},
				{
					input:  []byte(".1"),
					output: '1',
				},
				{
					input:  []byte(".1."),
					output: '1',
				},
				{
					input:  []byte("...1.."),
					output: '1',
				},
				{
					input:  []byte("x...1.."),
					output: 0,
					err:    common.NewParseError(0, "expected digit"),
				},
				{
					input:  []byte("...1..x"),
					output: '1',
				},
			},
		},
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	runTests(t, []test[bool]{
		{
			comb: EOF(),
			cases: []testCase[bool]{
				{
					input:  []byte{},
					output: true,
				},
				{
					input:  []byte("1"),
					output: false,
				},
				{
					input:  []byte("123"),
					output: false,
				},
			},
		},
	})
}

func TestCast(t *testing.T) {
	t.Parallel()

	runTests(t, []test[int]{
		{
			comb: Cast(
				Satisfy("test", true, common.Anything[byte]),
				func(x byte) (int, error) { return int(x), nil },
			),
			cases: []testCase[int]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "test"),
				},
				{
					input:  []byte{0},
					output: 0,
				},
				{
					input:  []byte{1},
					output: 1,
				},
				{
					input:  []byte{math.MaxUint8},
					output: 255,
				},
				{
					input:  []byte{math.MaxInt8},
					output: 127,
				},
			},
		},
		{
			comb: Cast(
				Any(),
				func(x byte) (int, error) { return -1, fmt.Errorf("test error") },
			),
			cases: []testCase[int]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, common.ErrEndOfFile.Error()),
				},
				{
					input:  []byte{0},
					output: 0,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{1},
					output: 0,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{math.MaxUint8},
					output: 0,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{math.MaxInt8},
					output: 0,
					err:    common.NewParseError(0, "test error"),
				},
			},
		},
	})
}
