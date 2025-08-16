package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestOr(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Or(
				"expected symbol 'a' or digit",
				Try(Eq("expected 'a'", 'a')),
				Range("expected digit", '0', '9'),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("3"),
					output: '3',
				},
				{
					input:  []byte("5"),
					output: '5',
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
			},
		},
		{
			comb: Or(
				"expected symbol 'a' or digit",
				Eq("expected 'a'", 'a'),
				Range("expected digit", '0', '9'),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("3"),
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
				{
					input:  []byte("5"),
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
				{
					input:  []byte("x3"),
					output: '3',
				},
				{
					input:  []byte("xz"),
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected symbol 'a' or digit"),
				},
			},
		},
	})
}

func TestAnd(t *testing.T) {
	runTestsSlice(t, []test[[]byte]{
		{
			comb: And(
				Eq("expected 'a'", 'a'),
				Eq("expected 'b'", 'b'),
				func(x, y byte) []byte { return []byte{x, y} },
			),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected 'a'"),
				},
				{
					input:  []byte("ab"),
					output: []byte{'a', 'b'},
				},
				{
					input:  []byte("abc"),
					output: []byte{'a', 'b'},
				},
				{
					input:  []byte("a"),
					output: nil,
					err:    common.NewParseError(1, "expected 'b'"),
				},
				{
					input:  []byte("ac"),
					output: nil,
					err:    common.NewParseError(1, "expected 'b'"),
				},
				{
					input:  []byte(".ab"),
					output: nil,
					err:    common.NewParseError(0, "expected 'a'"),
				},
			},
		},
	})
}
