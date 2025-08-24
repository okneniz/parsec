package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okneniz/parsec/common"
)

type (
	testCase[T any] struct {
		input  string
		output T
		err    error
	}

	test[T any] struct {
		comb  common.Combinator[rune, Position, T]
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

					result, err := Parse([]rune(testCase.input), test.comb)

					if testCase.err != nil {
						assert.Equal(t, err.Error(), testCase.err.Error())
					} else {
						assert.NoError(t, err)
					}

					assert.Equal(t, testCase.output, result)
				})
			}
		})
	}
}

func runTestsString[T comparable](t *testing.T, tests []test[[]T]) {
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

					result, err := ParseString(testCase.input, test.comb)

					if testCase.err != nil {
						assert.EqualError(t, err, testCase.err.Error())
					} else {
						assert.NoError(t, err)
					}

					assert.EqualValues(t, testCase.output, result)
				})
			}
		})
	}
}

func TestRange(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Range("expected char between 'a' and 'c'", 'a', 'c'),
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
						"expected char between 'a' and 'c'",
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
						"expected char between 'a' and 'c'",
					),
				},
				{
					input:  "ad",
					output: 'a',
				},
				{
					input:  "da",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char between 'a' and 'c'",
					),
				},
			},
		},
	})
}

func TestNotRange(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: NotRange("expected char not between 'a' and 'c'", 'a', 'c'),
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
						"expected char not between 'a' and 'c'",
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
						"expected char not between 'a' and 'c'",
					),
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
						"expected char not between 'a' and 'c'",
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
						"expected char not between 'a' and 'c'",
					),
				},
				{
					input:  "d",
					output: 'd',
				},
				{
					input:  "ad",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char not between 'a' and 'c'",
					),
				},
				{
					input:  "da",
					output: 'd',
				},
			},
		},
	})
}

func TestGt(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Gt("expected char greater than 'c'", 'c'),
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
						"expected char greater than 'c'",
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
						"expected char greater than 'c'",
					),
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
						"expected char greater than 'c'",
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
						"expected char greater than 'c'",
					),
				},
				{
					input:  "d",
					output: 'd',
				},
				{
					input:  "e",
					output: 'e',
				},
				{
					input:  "ad",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char greater than 'c'",
					),
				},
				{
					input:  "da",
					output: 'd',
				},
			},
		},
	})
}

func TestGte(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Gte("expected char greater than or equal 'c'", 'c'),
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
						"expected char greater than or equal 'c'",
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
						"expected char greater than or equal 'c'",
					),
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
						"expected char greater than or equal 'c'",
					),
				},
				{
					input:  "c",
					output: 'c',
				},
				{
					input:  "d",
					output: 'd',
				},
				{
					input:  "e",
					output: 'e',
				},
				{
					input:  "ad",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char greater than or equal 'c'",
					),
				},
				{
					input:  "da",
					output: 'd',
				},
			},
		},
	})
}

func TestLt(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Lt("expected char less than 'c'", 'c'),
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
						"expected char less than 'c'",
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
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char less than 'c'",
					),
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
						"expected char less than 'c'",
					),
				},
				{
					input:  "e",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char less than 'c'",
					),
				},
				{
					input:  "ad",
					output: 'a',
				},
				{
					input:  "da",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char less than 'c'",
					),
				},
			},
		},
	})
}

func TestLte(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Lte("expected char less than or equal 'c'", 'c'),
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
						"expected char less than or equal 'c'",
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
						"expected char less than or equal 'c'",
					),
				},
				{
					input:  "e",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char less than or equal 'c'",
					),
				},
				{
					input:  "ad",
					output: 'a',
				},
				{
					input:  "da",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected char less than or equal 'c'",
					),
				},
			},
		},
	})
}
