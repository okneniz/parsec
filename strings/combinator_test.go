package strings

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
	"unicode"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"
	"github.com/okneniz/parsec"
)

func TestSatisfy(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Satisfy[parsec.Stateless]("expected not 'c'", true, func(x rune) bool { return x != 'c' }),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Satisfy[parsec.Stateless]("error explanation", true, parsec.Nothing),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Any[parsec.Stateless](),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						parsec.ErrEndOfFile.Error(),
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
		comb := Any[parsec.Stateless]()

		ohsnap.Check(t, 1000, arbitrary.Rune(rnd, 0, unicode.MaxRune), func(r rune) bool {
			result, err := Parse(parsec.Stateless{}, []rune{r}, comb)
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
				Satisfy[parsec.Stateless](
					"error explanation",
					true,
					func(x rune) bool { return x <= 'b' },
				),
			),
			cases: []testCase[rune]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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

	notBrackets := Satisfy[parsec.Stateless]("test", true, func(x rune) bool {
		return x != ')' && x != '('
	})

	comb := Between(
		Eq[parsec.Stateless]("expected '('", '('),
		Some(0, "expected not ( or ) symbols", Try(notBrackets)),
		Eq[parsec.Stateless]("expected ')'", ')'),
	)

	runTestsString(t, []test[[]rune]{
		{
			comb: comb,
			cases: []testCase[[]rune]{
				{
					input:  "",
					output: nil,
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: EOF[parsec.Stateless](),
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
				Satisfy[parsec.Stateless]("test", true, parsec.Anything[rune]),
				func(x rune) (int, error) { return int(x), nil },
			),
			cases: []testCase[int]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
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
				Any[parsec.Stateless](),
				func(x rune) (int, error) { return -1, fmt.Errorf("test error") },
			),
			cases: []testCase[int]{
				{
					input:  "",
					output: 0,
					err: parsec.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						parsec.ErrEndOfFile.Error(),
					),
				},
				{
					input:  "something",
					output: -1,
					err: parsec.NewParseError(
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
