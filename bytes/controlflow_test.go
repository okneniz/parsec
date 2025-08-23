package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestConcat(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: Concat(
				6,
				Count(1, "expected one a", Eq("expected a", 'a')),
				Count(2, "expected double b", Eq("expected b", 'b')),
				Count(3, "expected three not z", NotEq("expected not z", 'z')),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected one a"),
				},
				{
					input:  []byte("abbcde"),
					output: []byte{'a', 'b', 'b', 'c', 'd', 'e'},
				},
				{
					input: []byte("x"),
					err:   common.NewParseError(0, "expected a"),
				},
				{
					input: []byte("ax"),
					err:   common.NewParseError(1, "expected b"),
				},
				{
					input: []byte("ab"),
					err:   common.NewParseError(1, "expected double b"),
				},
				{
					input: []byte("abb"),
					err:   common.NewParseError(3, "expected three not z"),
				},
				{
					input: []byte("abbc"),
					err:   common.NewParseError(3, "expected three not z"),
				},
				{
					input: []byte("abbcd"),
					err:   common.NewParseError(3, "expected three not z"),
				},
				{
					input:  []byte("abbcde"),
					output: []byte{'a', 'b', 'b', 'c', 'd', 'e'},
				},
				{
					input:  []byte("abbcdez"),
					output: []byte{'a', 'b', 'b', 'c', 'd', 'e'},
				},
				{
					input:  []byte("abbcdef"),
					output: []byte{'a', 'b', 'b', 'c', 'd', 'e'},
				},
			},
		},
	})
}

func TestSequence(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: Sequence(
				3,
				Eq("expected a", 'a'),
				Eq("expected b", 'b'),
				NotEq("expected not z", 'z'),
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected a"),
				},
				{
					input:  []byte("abc"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input:  []byte("abcd"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input: []byte(".abcd"),
					err:   common.NewParseError(0, "expected a"),
				},
				{
					input: []byte("a.bcd"),
					err:   common.NewParseError(1, "expected b"),
				},
				{
					input: []byte("abzcd"),
					err:   common.NewParseError(2, "expected not z"),
				},
				{
					input:  []byte("abcz"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input: []byte("ab"),
					err:   common.NewParseError(2, "expected not z"),
				},
				{
					input: []byte("a"),
					err:   common.NewParseError(1, "expected b"),
				},
				{
					input: []byte("z"),
					err:   common.NewParseError(0, "expected a"),
				},
			},
		},
	})
}

func TestChoice(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Choice(
				Try(Eq("expected a", 'a')),
				Try(Eq("expected b", 'b')),
				Eq("expected c", 'c'),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected a"),
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
					output: 'c',
				},
				{
					input:  []byte("x"),
					output: 0,
					err:    common.NewParseError(0, "expected a"),
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
