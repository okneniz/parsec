package strings

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestSatisfy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Satisfy(true, func(x rune) bool { return x != 'c' })

		result, err := ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, 'a')

		result, err = ParseString("b", comb)
		Check(t, err)
		AssertEq(t, result, 'b')

		result, err = ParseString("c", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Satisfy(true, func(x rune) bool { return false })

		result, err := ParseString("", comb)
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
			b := rune(source.Intn(math.MaxUint8 + 1))

			result, err := Parse([]rune{rune(b)}, comb)
			Check(t, err)
			AssertEq(t, result, rune(b))
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Any()

		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Try(Satisfy(true, func(x rune) bool { return x <= 'b' }))

		// TODO : recover position assertions

		buf := Buffer([]rune("abcd"))
		// AssertEq(t, buf.Position(), 0)

		result, err := p.Parse[rune, Position, rune](buf, comb)
		Check(t, err)
		AssertEq(t, result, 'a')
		// AssertEq(t, buf.Position(), 1)

		result, err = p.Parse[rune, Position, rune](buf, comb)
		Check(t, err)
		AssertEq(t, result, 'b')
		// AssertEq(t, buf.Position(), 2)

		result, err = p.Parse[rune, Position, rune](buf, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
		// AssertEq(t, buf.Position(), 2)

		result, err = p.Parse[rune, Position, rune](buf, comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
		// AssertEq(t, buf.Position(), 2)
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		notBrackets := Satisfy(true, func(x rune) bool {
			return !(x == ')' || x == '(')
		})

		comb := Between(
			Eq('('),
			Some(0, Try(notBrackets)),
			Eq(')'),
		)

		result, err := ParseString("(abc)", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("(abc)def", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("(abc))", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("(ab)", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b'})

		result, err = ParseString("x(abc)def", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("()", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("(()", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("((1))", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("(abc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("(abc(", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("((abc)", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestPadded(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Padded(
			Eq('.'),
			Range('0', '9'),
		)

		result, err := ParseString("1", comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = ParseString(".1", comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = ParseString("...1", comb)
		Check(t, err)
		AssertEq(t, result, '1')

		result, err = ParseString("..1...", comb)
		Check(t, err)
		AssertEq(t, result, '1')
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("abcd", EOF())
		Check(t, err)
		AssertEq(t, result, false)
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("", EOF())
		Check(t, err)
		AssertEq(t, result, true)
	})
}

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Cast(
			Satisfy(true, p.Anything[rune]),
			func(x rune) (int, error) { return int(x), nil },
		)

		result, err := ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, 97)

		result, err = ParseString("b", comb)
		Check(t, err)
		AssertEq(t, result, 98)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Cast(
			Any(),
			func(x rune) (int, error) { return -1, fmt.Errorf("test error") },
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}
