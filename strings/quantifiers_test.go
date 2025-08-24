package strings

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestMany(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: Many(0, Eq("expected a", 'a')),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
				},
				{
					input:  "a",
					output: []rune{'a'},
				},
				{
					input:  "aaa",
					output: []rune{'a', 'a', 'a'},
				},
				{
					input:  "aaab",
					output: []rune{'a', 'a', 'a'},
				},
				{
					input:  "aaa.aa",
					output: []rune{'a', 'a', 'a'},
				},
				{
					input:  ".aaa",
					output: []rune{},
				},
			},
		},
	})
}

func TestSome(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: Some(0, "expected at least one 'a'", Eq("expected 'a'", 'a')),
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
						"expected at least one 'a'",
					),
				},
				{
					input:  "a",
					output: []rune{'a'},
				},
				{
					input:  "aaa",
					output: []rune{'a', 'a', 'a'},
				},
				{
					input:  "aa.",
					output: []rune{'a', 'a'},
				},
				{
					input:  "aa.aaa",
					output: []rune{'a', 'a'},
				},
				{
					input:  ".aa",
					output: nil,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected at least one 'a'",
					),
				},
			},
		},
	})
}

func TestOptional(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Optional(Eq("expected a", 'a'), 123),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 123,
				},
				{
					input:  "a",
					output: 'a',
				},
				{
					input:  "aa",
					output: 'a',
				},
				{
					input:  "xa",
					output: 123,
				},
				{
					input:  "ax",
					output: 'a',
				},
			},
		},
		{
			comb: Optional(Satisfy("never match", true, common.Nothing[rune]), 'x'),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 'x',
				},
				{
					input:  "a",
					output: 'x',
				},
				{
					input:  "aa",
					output: 'x',
				},
				{
					input:  "za",
					output: 'x',
				},
				{
					input:  "az",
					output: 'x',
				},
			},
		},
	})
}

func lTestCount(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: Count(2, "expected 'aa'", Eq("expected 'a'", 'a')),
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: []rune{},
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected 'aa'",
					),
				},
				{
					input:  "aa",
					output: []rune{'a', 'a'},
				},
				{
					input:  "aaa",
					output: []rune{'a', 'a'},
				},
				{
					input:  "aa.",
					output: []rune{'a', 'a'},
				},
				{
					input:  ".aa",
					output: []rune{},
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
					input:  "a.",
					output: []rune{},
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected 'a'",
					),
				},
			},
		},
	})
}
