package strings

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestConcat(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: Concat(
				6,
				Count(1, "expected one a", Eq("expected a", 'a')),
				Count(2, "expected double b", Eq("expected b", 'b')),
				Count(3, "expected three not z", NotEq("expected not z", 'z')),
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
						"expected one a",
					),
				},
				{
					input:  "abbcde",
					output: []rune("abbcde"),
				},
				{
					input: "x",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected one a",
					),
				},
				{
					input: "ax",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected double b",
					),
				},
				{
					input: "ab",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected double b",
					),
				},
				{
					input: "abb",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 3,
							index:  3,
						},
						"expected three not z",
					),
				},
				{
					input: "abbc",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 3,
							index:  3,
						},
						"expected three not z",
					),
				},
				{
					input: "abbcd",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 3,
							index:  3,
						},
						"expected three not z",
					),
				},
				{
					input:  "abbcde",
					output: []rune("abbcde"),
				},
				{
					input:  "abbcdez",
					output: []rune("abbcde"),
				},
				{
					input:  "abbcdef",
					output: []rune("abbcde"),
				},
			},
		},
	})
}

func TestSequence(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: Sequence(
				3,
				Eq("expected a", 'a'),
				Eq("expected b", 'b'),
				NotEq("expected not z", 'z'),
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
						"expected a",
					),
				},
				{
					input:  "abc",
					output: []rune("abc"),
				},
				{
					input:  "abcd",
					output: []rune("abc"),
				},
				{
					input: ".abcd",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected a",
					),
				},
				{
					input: "a.bcd",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected b",
					),
				},
				{
					input: "abzcd",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected not z",
					),
				},
				{
					input:  "abcz",
					output: []rune("abc"),
				},
				{
					input: "ab",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected not z",
					),
				},
				{
					input: "a",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected b",
					),
				},
				{
					input: "z",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected a",
					),
				},
			},
		},
	})
}

func TestChoice(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Choice(
				"expected a, b or c",
				Try(Eq("expected a", 'a')),
				Try(Eq("expected b", 'b')),
				Eq("expected c", 'c'),
			),
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
						"expected a, b or c",
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
					input:  "x",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected a, b or c",
					),
				},
			},
		},
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("skip not optional", func(t *testing.T) {
		runTests(t, []test[rune]{
			{
				comb: Skip(
					Eq("expected 'a'", 'a'),
					Eq("expected 'b'", 'b'),
				),
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
							"expected 'a'",
						),
					},
					{
						input:  "abc",
						output: 'b',
					},
					{
						input:  "b",
						output: 0,
						err: common.NewParseError(
							Position{
								line:   0,
								column: 0,
								index:  0,
							},
							"expected 'a'",
						),
					},
					{
						input: "bbb",
						err: common.NewParseError(
							Position{
								line:   0,
								column: 0,
								index:  0,
							},
							"expected 'a'",
						),
					},
					{
						input:  "ac",
						output: 0,
						err: common.NewParseError(
							Position{
								line:   0,
								column: 1,
								index:  1,
							},
							"expected 'b'",
						),
					},
				},
			},
		})
	})

	t.Run("skip optional", func(t *testing.T) {
		runTests(t, []test[rune]{
			{
				comb: Skip(
					Optional(Try(Eq("expected 'a'", 'a')), 'x'),
					Eq("expected 'b'", 'b'),
				),
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
							"expected 'b'",
						),
					},
					{
						input:  "abc",
						output: 'b',
					},
					{
						input:  "b",
						output: 'b',
					},
					{
						input:  "bbb",
						output: 'b',
					},
					{
						input:  "ac",
						output: 0,
						err: common.NewParseError(
							Position{
								line:   0,
								column: 1,
								index:  1,
							},
							"expected 'b'",
						),
					},
				},
			},
		})
	})
}

func TestSkipMany(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: common.SkipMany(
				NoneOf("expected not a, b or c", 'a', 'b', 'c'),
				SequenceOf("expected abc", 'a', 'b', 'c'),
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
						"expected abc",
					),
				},
				{
					input:  "abc",
					output: []rune("abc"),
				},
				{
					input:  "ab",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected abc",
					),
				},
				{
					input:  "xab",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected abc",
					),
				},
				{
					input:  "123abc",
					output: []rune("abc"),
				},
				{
					input:  "bcabc",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected abc",
					),
				},
				{
					input:  "abcabc",
					output: []rune("abc"),
				},
				{
					input:  "123abc123",
					output: []rune("abc"),
				},
				{
					input:  "123abcabc",
					output: []rune("abc"),
				},
				{
					input:  "123",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 3,
							index:  3,
						},
						"expected abc",
					),
				},
			},
		},
	})
}

func TestSkipAfter(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: SkipAfter(
				Eq("expected 'b'", 'b'),
				Eq("expected 'a'", 'a'),
			),
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
						"expected 'a'",
					),
				},
				{
					input:  "abc",
					output: 'a',
				},
				{
					input:  "ab",
					output: 'a',
				},
				{
					input: "a",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected 'b'",
					),
				},
				{
					input: "ac",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected 'b'",
					),
				},
				{
					input: "b",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a'",
					),
				},
				{
					input: "bc",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a'",
					),
				},
				{
					input: "bb",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'a'",
					),
				},
			},
		},
	})
}

func TestPadded(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Padded(
				Eq("expected dot", '.'),
				Range("expected digit", '0', '9'),
			),
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
						"expected digit",
					),
				},
				{
					input:  "1",
					output: '1',
				},
				{
					input:  ".1",
					output: '1',
				},
				{
					input:  ".1.",
					output: '1',
				},
				{
					input:  "...1..",
					output: '1',
				},
				{
					input:  "x...1..",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected digit",
					),
				},
				{
					input:  "...1..x",
					output: '1',
				},
			},
		},
	})
}
