package bytes

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestSatisfy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		c := byte('c')

		comb := Satisfy("test", true, func(x byte) bool { return x != c })

		result, err := Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, 'a')

		result, err = Parse([]byte("b"), comb)
		Check(t, err)
		AssertEq(t, result, 'b')

		result, err = Parse([]byte("c"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Satisfy("test", true, func(x byte) bool { return false })

		result, err := Parse([]byte{}, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		source := rand.New(rand.NewSource(time.Now().UnixNano()))
		comb := Any()

		for i := 0; i < 10000; i++ {
			b := byte(source.Intn(math.MaxUint8 + 1))

			result, err := Parse([]byte{b}, comb)
			Check(t, err)
			AssertEq(t, result, b)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Any()

		result, err := Parse([]byte{}, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Try(
			Satisfy("test", true, func(x byte) bool { return x <= byte('b') }),
		)

		buf := Buffer([]byte("abcd"))
		AssertEq(t, buf.Position(), 0)

		result, err := common.Parse[byte, int, byte](buf, comb)
		Check(t, err)
		AssertEq(t, result, byte('a'))
		AssertEq(t, buf.Position(), 1)

		result, err = common.Parse[byte, int, byte](buf, comb)
		Check(t, err)
		AssertEq(t, result, byte('b'))
		AssertEq(t, buf.Position(), 2)

		result, err = common.Parse[byte, int, byte](buf, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
		AssertEq(t, buf.Position(), 2)

		result, err = common.Parse[byte, int, byte](buf, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
		AssertEq(t, buf.Position(), 2)
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		notBrackets := Satisfy("test", true, func(x byte) bool {
			return !(x == byte(')') || x == byte('('))
		})

		comb := Between(
			Eq("expected '('", '('),
			Some(0, "expected not ( or ) symbols", Try(notBrackets)),
			Eq("expected ')'", ')'),
		)

		result, err := Parse([]byte("(abc)"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("(abc)def"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("(abc))"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("(ab)"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b'})

		result, err = Parse([]byte("x(abc)def"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("()"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("(()"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("((1))"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("(abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("(abc("), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("((abc)"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Skip(
			Optional(Eq("expected 'a'", 'a'), 0),
			Eq("expected 'b'", 'b'),
		)

		result, err := Parse([]byte("abc"), comb)
		Check(t, err)
		AssertEq(t, result, byte('b'))

		result, err = Parse([]byte("cba"), comb)
		Check(t, err)
		AssertEq(t, result, byte('b'))
	})

	t.Run("case 2", func(t *testing.T) {
		phrase := SequenceOf("expected abc", 'a', 'b', 'c')
		comb := common.SkipMany(
			NoneOf("expected not abc", 'a', 'b', 'c'),
			phrase,
		)

		result, err := Parse([]byte("abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("abc123"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("123abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("123abc123"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Skip(
			NotEq("expected not 'a'", 'a'),
			Eq("expected 'a'", 'a'),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestSkipAfter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SkipAfter(
			Eq("expected 'b'", 'b'),
			Eq("expected 'a'", 'a'),
		)

		result, err := Parse([]byte("abc"), comb)
		Check(t, err)
		AssertEq(t, result, byte('a'))

		result, err = Parse([]byte("ab"), comb)
		Check(t, err)
		AssertEq(t, result, byte('a'))

		result, err = Parse([]byte("a"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SkipAfter(
			Eq("expected 'b'", 'b'),
			Satisfy("test", true, common.Nothing[byte]),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SkipAfter(
			Satisfy("test", true, common.Nothing[byte]),
			Eq("expected 'a'", 'a'),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestPadded(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Padded(
			Eq("expected dot", '.'),
			Range("expected digit between 0 and 9", '0', '9'),
		)

		result, err := Parse([]byte("1"), comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = Parse([]byte(".1"), comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = Parse([]byte("...1"), comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = Parse([]byte("..1..."), comb)
		Check(t, err)
		AssertEq(t, result, '1')
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, err := Parse([]byte("abcd"), EOF())
		Check(t, err)
		AssertEq(t, result, false)
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := Parse([]byte(""), EOF())
		Check(t, err)
		AssertEq(t, result, true)
	})
}

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Cast(
			Satisfy("test", true, common.Anything[byte]),
			func(x byte) (int, error) { return int(x), nil },
		)

		result, err := Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, 97)

		result, err = Parse([]byte("b"), comb)
		Check(t, err)
		AssertEq(t, result, 98)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Cast(
			Any(),
			func(x byte) (int, error) { return -1, fmt.Errorf("test error") },
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}
