package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestEq(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Eq("expected 'c'", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected 'c'"),
				},
				{
					input:  []byte("a"),
					output: 0,
					err:    common.NewParseError(0, "expected 'c'"),
				},
				{
					input:  []byte("c"),
					output: 'c',
				},
				{
					input:  []byte("ca"),
					output: 'c',
				},
				{
					input:  []byte("ac"),
					output: 0,
					err:    common.NewParseError(0, "expected 'c'"),
				},
			},
		},
	})
}

func TestNotEq(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: NotEq("expected not c", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected not c"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected not c"),
				},
				{
					input:  []byte("ca"),
					output: 0,
					err:    common.NewParseError(0, "expected not c"),
				},
				{
					input:  []byte("ac"),
					output: 'a',
				},
			},
		},
	})
}

func TestOneOf(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: OneOf("expected 'a', 'b' or 'c'", 'a', 'b', 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected 'a', 'b' or 'c'"),
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
					input:  []byte("d"),
					output: 0,
					err:    common.NewParseError(0, "expected 'a', 'b' or 'c'"),
				},
				{
					input:  []byte("ca"),
					output: 'c',
				},
				{
					input:  []byte("ac"),
					output: 'a',
				},
				{
					input:  []byte("bb"),
					output: 'b',
				},
				{
					input: []byte("fa"),
					err:   common.NewParseError(0, "expected 'a', 'b' or 'c'"),
				},
			},
		},
	})
}

func TestSequenceOf(t *testing.T) {
	runTestsSlice(t, []test[[]byte]{
		{
			comb: SequenceOf("expected foo", 'f', 'o', 'o'),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
				{
					input:  []byte{' '},
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
				{
					input:  []byte("f"),
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
				{
					input:  []byte("fo"),
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
				{
					input:  []byte("foo"),
					output: []byte{'f', 'o', 'o'},
				},
				{
					input:  []byte("foo."),
					output: []byte{'f', 'o', 'o'},
				},
				{
					input:  []byte(".foo"),
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
				{
					input:  []byte("foobar"),
					output: []byte{'f', 'o', 'o'},
				},
				{
					input:  []byte("barfoo"),
					output: nil,
					err:    common.NewParseError(0, "expected foo"),
				},
			},
		},
	})
}

func TestMap(t *testing.T) {
	runTestsSlice(t, []test[[]string]{
		{
			comb: Some(
				1,
				"expected at least one 0, 1 or 2",
				common.SkipMany(
					NoneOf("skip not 0, 1 or 2", 0, 1, 2),
					Map(
						"expected 0, 1 or 2",
						map[byte]string{
							0: "foo",
							1: "bar",
							2: "baz",
						},
						Any(),
					),
				),
			),
			cases: []testCase[[]string]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected at least one 0, 1 or 2"),
				},
				{
					input:  []byte{0},
					output: []string{"foo"},
				},
				{
					input:  []byte{0, 1},
					output: []string{"foo", "bar"},
				},
				{
					input:  []byte{0, 1, 2},
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  []byte{0, 1, 2, 3},
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  []byte{10, 0, 1, 2, 3, 4, 5, 1, 2, 3, 0, 1},
					output: []string{"foo", "bar", "baz", "bar", "baz", "foo", "bar"},
				},
				{
					input:  []byte{5, 6, 7, 8, 9, 10, 11},
					output: nil,
					err:    common.NewParseError(0, "expected at least one 0, 1 or 2"),
				},
			},
		},
	})
}
