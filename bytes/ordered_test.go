package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
)

func TestRange(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Range("expected any byte between 'a' and 'c'", 'a', 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected any byte between 'a' and 'c'"),
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
					err:    common.NewParseError(0, "expected any byte between 'a' and 'c'"),
				},
				{
					input:  []byte("da"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte between 'a' and 'c'"),
				},
				{
					input:  []byte("ad"),
					output: 'a',
				},
			},
		},
	})
}

func TestNotRange(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: NotRange("expected any byte not between 'a' and 'c'", 'a', 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected any byte not between 'a' and 'c'"),
				},
				{
					input:  []byte("a"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte not between 'a' and 'c'"),
				},
				{
					input:  []byte("b"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte not between 'a' and 'c'"),
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte not between 'a' and 'c'"),
				},
				{
					input:  []byte("d"),
					output: 'd',
				},
				{
					input:  []byte("da"),
					output: 'd',
				},
				{
					input: []byte("ad"),
					err:   common.NewParseError(0, "expected any byte not between 'a' and 'c'"),
				},
			},
		},
	})
}

func TestGt(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Gt("expected any byte greater than 'c'", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than 'c'"),
				},
				{
					input:  []byte("a"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than 'c'"),
				},
				{
					input:  []byte("b"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than 'c'"),
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than 'c'"),
				},
				{
					input:  []byte("d"),
					output: 'd',
				},
				{
					input:  []byte("da"),
					output: 'd',
				},
				{
					input:  []byte("e"),
					output: 'e',
				},
				{
					input: []byte("ad"),
					err:   common.NewParseError(0, "expected any byte greater than 'c'"),
				},
			},
		},
	})
}

func TestGte(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Gte("expected any byte greater than or equal 'c'", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than or equal 'c'"),
				},
				{
					input:  []byte("a"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than or equal 'c'"),
				},
				{
					input:  []byte("b"),
					output: 0,
					err:    common.NewParseError(0, "expected any byte greater than or equal 'c'"),
				},
				{
					input:  []byte("c"),
					output: 'c',
				},
				{
					input:  []byte("d"),
					output: 'd',
				},
				{
					input:  []byte("da"),
					output: 'd',
				},
				{
					input:  []byte("e"),
					output: 'e',
				},
				{
					input: []byte("ad"),
					err:   common.NewParseError(0, "expected any byte greater than or equal 'c'"),
				},
			},
		},
	})
}

func TestLt(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Lt("expected byte less than 'c'", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected byte less than 'c'"),
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
					output: 0,
					err:    common.NewParseError(0, "expected byte less than 'c'"),
				},
				{
					input:  []byte("d"),
					output: 0,
					err:    common.NewParseError(0, "expected byte less than 'c'"),
				},
				{
					input:  []byte("da"),
					output: 0,
					err:    common.NewParseError(0, "expected byte less than 'c'"),
				},
			},
		},
	})
}

func TestLte(t *testing.T) {
	runTests(t, []test[byte]{
		{
			comb: Lte("expected byte less than or equal 'c'", 'c'),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected byte less than or equal 'c'"),
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
					err:    common.NewParseError(0, "expected byte less than or equal 'c'"),
				},
				{
					input:  []byte("da"),
					output: 0,
					err:    common.NewParseError(0, "expected byte less than or equal 'c'"),
				},
			},
		},
	})
}
