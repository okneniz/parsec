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
