package strings

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestEq(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Eq("expected 'c'", 'c'),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'c'",
					),
				},
				{
					input:  "a",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'c'",
					),
				},
				{
					input:  "c",
					output: 'c',
				},
				{
					input:  "ca",
					output: 'c',
				},
				{
					input:  "ac",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'c'",
					),
				},
			},
		},
	})
}

func TestNotEq(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: NotEq("expected not c", 'c'),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected not c",
					),
				},
				{
					input:  "a",
					output: 'a',
				},
				{
					input:  "c",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected not c",
					),
				},
				{
					input:  "ca",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected not c",
					),
				},
				{
					input:  "ac",
					output: 'a',
				},
			},
		},
	})
}

func TestOneOf(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: OneOf("expected 'a', 'b' or 'c'", 'a', 'b', 'c'),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a', 'b' or 'c'",
					),
				},
				{
					input:  "a",
					output: 'a',
				},
				{
					input:  "b",
					output: 'b',
				},
				{
					input:  "c",
					output: 'c',
				},
				{
					input:  "d",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a', 'b' or 'c'",
					),
				},
				{
					input:  "ca",
					output: 'c',
				},
				{
					input:  "ac",
					output: 'a',
				},
				{
					input:  "bb",
					output: 'b',
				},
				{
					input: "fa",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a', 'b' or 'c'",
					),
				},
			},
		},
	})
}

func TestSequenceOf(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: SequenceOf("expected foo", 'f', 'o', 'o'),
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
						"expected foo",
					),
				},
				{
					input:  " ",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "f",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "fo",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "foo",
					output: []rune{'f', 'o', 'o'},
				},
				{
					input:  "foo.",
					output: []rune{'f', 'o', 'o'},
				},
				{
					input:  ".foo",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "foobar",
					output: []rune{'f', 'o', 'o'},
				},
				{
					input:  "barfoo",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
			},
		},
	})
}

func TestMap(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]string]{
		{
			comb: Some(
				1,
				"expected at least one a, b or c",
				common.SkipMany(
					NoneOf("skip not a, b or c", 'a', 'b', 'c'),
					Map(
						"expected a, b or c",
						map[rune]string{
							'a': "foo",
							'b': "bar",
							'c': "baz",
						},
						Any(),
					),
				),
			),
			cases: []testCase[[]string]{
				{
					input:  "",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one a, b or c",
					),
				},
				{
					input:  "a",
					output: []string{"foo"},
				},
				{
					input:  "ab",
					output: []string{"foo", "bar"},
				},
				{
					input:  "abc",
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  "abcd",
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  "abcbcabzx",
					output: []string{"foo", "bar", "baz", "bar", "baz", "foo", "bar"},
				},
				{
					input:  "xyzsert",
					output: nil,

					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one a, b or c",
					),
				},
			},
		},
	})
}

func TestMapStrings(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]string]{
		{
			comb: Some(
				1,
				"sequence of keys",
				SkipMany(
					NoneOf(
						"none of 'a', 'b' or 'c'",
						'a', 'b', 'c',
					),
					MapStrings(
						"expect 'a', 'b' or 'c'",
						map[string]string{
							"a": "foo",
							"b": "bar",
							"c": "baz",
						},
					),
				),
			),
			cases: []testCase[[]string]{
				{
					input:  "",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"sequence of keys",
					),
				},
				{
					input:  "a",
					output: []string{"foo"},
				},
				{
					input:  "ab",
					output: []string{"foo", "bar"},
				},
				{
					input:  "abc",
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  "abcd",
					output: []string{"foo", "bar", "baz"},
				},
				{
					input:  "abcbcabzx",
					output: []string{"foo", "bar", "baz", "bar", "baz", "foo", "bar"},
				},
				{
					input:  "..a//b++c**d,,e--a",
					output: []string{"foo", "bar", "baz", "foo"},
				},
				{
					input:  "xyzsert",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"sequence of keys",
					),
				},
			},
		},
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	runTests(t, []test[string]{
		{
			comb: String("expected foo", "foo"),
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
						"expected foo",
					),
				},
				{
					input:  "foo",
					output: "foo",
				},
				{
					input:  "foobar",
					output: "foo",
				},
				{
					input:  "fo",
					output: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "f",
					output: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
				{
					input:  "bar",
					output: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected foo",
					),
				},
			},
		},
	})
}
