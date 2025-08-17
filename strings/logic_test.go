package strings

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestOr(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Or(
				"expected symbol 'a' or digit",
				Try(Eq("expected 'a'", 'a')),
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
						"expected symbol 'a' or digit",
					),
				},
				{
					input:  "a",
					output: 'a',
				},
				{
					input:  "3",
					output: '3',
				},
				{
					input:  "5",
					output: '5',
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
						"expected symbol 'a' or digit",
					),
				},
			},
		},
		{
			comb: Or(
				"expected symbol 'a' or digit",
				Eq("expected 'a'", 'a'),
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
						"expected symbol 'a' or digit",
					),
				},
				{
					input:  "a",
					output: 'a',
				},
				{
					input:  "3",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected symbol 'a' or digit",
					),
				},
				{
					input:  "5",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected symbol 'a' or digit",
					),
				},
				{
					input:  "x3",
					output: '3',
				},
				{
					input:  "xz",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected symbol 'a' or digit",
					),
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
						"expected symbol 'a' or digit",
					),
				},
			},
		},
	})
}

func TestAnd(t *testing.T) {
	t.Parallel()

	runTestsString(t, []test[[]rune]{
		{
			comb: And(
				Eq("expected 'a'", 'a'),
				Eq("expected 'b'", 'b'),
				func(x, y rune) []rune { return []rune{x, y} },
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
						"expected 'a'",
					),
				},
				{
					input:  "ab",
					output: []rune{'a', 'b'},
				},
				{
					input:  "abc",
					output: []rune{'a', 'b'},
				},
				{
					input:  "a",
					output: nil,
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
					input:  "ac",
					output: nil,
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
					input:  ".ab",
					output: nil,
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
