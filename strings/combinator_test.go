package strings

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec/common"
)

func TestSatisfy(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Satisfy("expected not 'c'", true, func(x rune) bool { return x != 'c' }),
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
						"expected not 'c'",
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
						"expected not 'c'",
					),
				},
			},
		},
		{
			comb: Satisfy("error explanation", true, func(x rune) bool { return false }),
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
						"error explanation",
					),
				},
				{
					input:  "abc",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"error explanation",
					),
				},
			},
		},
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Any(),
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
						common.ErrEndOfFile.Error(),
					),
				},
				{
					input:  "a",
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

		ohsnap.Check(t, 1000, ohsnap.NewBuilder(rnd).Rune(), func(r rune) bool {
			result, err := Parse([]rune{r}, comb)
			if err != nil {
				t.Logf("input: %v", []rune{r})
				t.Logf("output: %v", result)
				t.Error(err)
				return false
			}

			return result == r
		})
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Try(
				Satisfy(
					"error explanation",
					true,
					func(x rune) bool { return x <= 'b' },
				),
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
						"error explanation",
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
					input: "c",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"error explanation",
					),
				},
			},
		},
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	notBrackets := Satisfy("test", true, func(x rune) bool {
		return !(x == ')' || x == '(')
	})

	comb := Between(
		Eq("expected '('", '('),
		Some(0, "expected not ( or ) symbols", Try(notBrackets)),
		Eq("expected ')'", ')'),
	)

	runTestsString(t, []test[[]rune]{
		{
			comb: comb,
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
						"expected '('",
					),
				},
				{
					input:  "(abc)",
					output: []rune{'a', 'b', 'c'},
				},
				{
					input: "abc",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '('",
					),
				},
				{
					input: "(abc",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 4,
							index:  4,
						},
						"expected ')'",
					),
				},
				{
					input: " (abc) ",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '('",
					),
				},
				{
					input: "((abc))",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected not ( or ) symbols",
					),
				},
				{
					input: "()",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected not ( or ) symbols",
					),
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
					input:  "",
					output: true,
				},
				{
					input:  "1",
					output: false,
				},
				{
					input:  "123",
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
				Satisfy("test", true, common.Anything[rune]),
				func(x rune) (int, error) { return int(x), nil },
			),
			cases: []testCase[int]{
				{
					input:  "",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"test",
					),
				},
				{
					input:  string([]rune{0}),
					output: 0,
				},
				{
					input:  "a",
					output: 97,
				},
			},
		},
		{
			comb: Cast(
				Any(),
				func(x rune) (int, error) { return -1, fmt.Errorf("test error") },
			),
			cases: []testCase[int]{
				{
					input:  "",
					output: 0,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						common.ErrEndOfFile.Error(),
					),
				},
				{
					input:  "something",
					output: -1,
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"test error",
					),
				},
			},
		},
	})
}
