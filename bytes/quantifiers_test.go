package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestMany(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: Many(0, Eq("expected a", 'a')),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
				},
				{
					input:  []byte("a"),
					output: []byte{'a'},
				},
				{
					input:  []byte("aaa"),
					output: []byte{'a', 'a', 'a'},
				},
				{
					input:  []byte("aaab"),
					output: []byte{'a', 'a', 'a'},
				},
				{
					input:  []byte("aaa.aa"),
					output: []byte{'a', 'a', 'a'},
				},
				{
					input:  []byte(".aaa"),
					output: nil,
				},
			},
		},
	})
}

func TestSome(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: Some(0, "expected at least one 'a'", Eq("expected 'a'", 'a')),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected at least one 'a'"),
				},
				{
					input:  []byte("a"),
					output: []byte{'a'},
				},
				{
					input:  []byte("aaa"),
					output: []byte{'a', 'a', 'a'},
				},
				{
					input:  []byte("aa."),
					output: []byte{'a', 'a'},
				},
				{
					input:  []byte("aa.aaa"),
					output: []byte{'a', 'a'},
				},
				{
					input:  []byte(".aa"),
					output: nil,
					err:    common.NewParseError(0, "expected at least one 'a'"),
				},
			},
		},
	})
}

func TestOptional(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Optional(Eq("expected a", 'a'), 123),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 123,
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("aa"),
					output: 'a',
				},
				{
					input:  []byte("xa"),
					output: 123,
				},
				{
					input:  []byte("ax"),
					output: 'a',
				},
			},
		},
		{
			comb: Optional(Satisfy("never match", true, common.Nothing[byte]), 'x'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 'x',
				},
				{
					input:  []byte("a"),
					output: 'x',
				},
				{
					input:  []byte("aa"),
					output: 'x',
				},
				{
					input:  []byte("za"),
					output: 'x',
				},
				{
					input:  []byte("az"),
					output: 'x',
				},
			},
		},
	})
}

func TestCount(t *testing.T) {
	t.Parallel()

	runTestsSlice(t, []test[[]byte]{
		{
			comb: Count(2, "expected 'aa'", Eq("expected 'a'", 'a')),
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected 'aa'"),
				},
				{
					input:  []byte("aa"),
					output: []byte{'a', 'a'},
				},
				{
					input:  []byte("aaa"),
					output: []byte{'a', 'a'},
				},
				{
					input:  []byte("aa."),
					output: []byte{'a', 'a'},
				},
				{
					input:  []byte(".aa"),
					output: nil,
					err:    common.NewParseError(0, "expected 'a'"),
				},
				{
					input:  []byte("a."),
					output: nil,
					err:    common.NewParseError(1, "expected 'a'"),
				},
			},
		},
	})
}
