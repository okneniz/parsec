package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestParens(t *testing.T) {
	t.Parallel()

	comb := Parens(Some(1, Try(IsDigit())))

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("(123)", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("123"))
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("(123", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("123)", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 4", func(t *testing.T) {
		result, err := ParseString("(asd)", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 5", func(t *testing.T) {
		result, err := ParseString("()", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 6", func(t *testing.T) {
		result, err := ParseString("(", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 7", func(t *testing.T) {
		result, err := ParseString(")", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 8", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestBraces(t *testing.T) {
	t.Parallel()

	comb := Braces(Some(1, Try(IsDigit())))

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("{123}", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("123"))
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("{123", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("123}", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 4", func(t *testing.T) {
		result, err := ParseString("{asd}", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 5", func(t *testing.T) {
		result, err := ParseString("{}", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 6", func(t *testing.T) {
		result, err := ParseString("{", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 7", func(t *testing.T) {
		result, err := ParseString("}", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 8", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestAngles(t *testing.T) {
	t.Parallel()

	comb := Angles(Some(1, Try(IsDigit())))

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("<123>", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("123"))
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("<123", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("123>", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 4", func(t *testing.T) {
		result, err := ParseString("<asd>", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 5", func(t *testing.T) {
		result, err := ParseString("<>", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 6", func(t *testing.T) {
		result, err := ParseString("<", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 7", func(t *testing.T) {
		result, err := ParseString(">", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 8", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSquares(t *testing.T) {
	t.Parallel()

	comb := Squares(Some(1, Try(IsDigit())))

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("[123]", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("123"))
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("[123", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("123]", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 4", func(t *testing.T) {
		result, err := ParseString("[asd]", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 5", func(t *testing.T) {
		result, err := ParseString("[]", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 6", func(t *testing.T) {
		result, err := ParseString("[", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 7", func(t *testing.T) {
		result, err := ParseString("]", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 8", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSemi(t *testing.T) {
	t.Parallel()

	comb := SepBy1(3, NotEq(';'), Semi())

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("1;2;3", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1', '2', '3'})
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("1", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1'})
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestComma(t *testing.T) {
	t.Parallel()

	comb := SepBy1(3, NotEq(','), Comma())

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("1,2,3", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1', '2', '3'})
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("1", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1'})
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestColon(t *testing.T) {
	t.Parallel()

	comb := SepBy1(3, NotEq(':'), Colon())

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("1:2:3", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1', '2', '3'})
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("1", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1'})
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestDot(t *testing.T) {
	t.Parallel()

	comb := SepBy1(3, NotEq('.'), Dot())

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseString("1.2.3", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1', '2', '3'})
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseString("1", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'1'})
	})

	t.Run("case 3", func(t *testing.T) {
		result, err := ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestUnsigned(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		for s, i := range map[string]uint{
			"0":      0,
			"1":      1,
			"2":      2,
			"3":      3,
			"4":      4,
			"5":      5,
			"6":      6,
			"7":      7,
			"8":      8,
			"9":      9,
			"10":     10,
			"500":    500,
			"100500": 100500,
		} {
			result, err := ParseString(s, Unsigned[uint]())
			Check(t, err)
			AssertEq(t, result, i)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, s := range []string{
			"asdasd10",
			" 50OIUO",
			"(10)(*(0))",
			"_10asd",
		} {
			result, err := ParseString(s, Unsigned[uint]())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestUnsignedN(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		for s, i := range map[string]uint{
			"0":      0,
			"1":      1,
			"2":      2,
			"3":      3,
			"4":      4,
			"5":      5,
			"6":      6,
			"7":      7,
			"8":      8,
			"9":      9,
			"10":     1,
			"500":    5,
			"100500": 1,
		} {
			result, err := ParseString(s, UnsignedN[uint](1))
			Check(t, err)
			AssertEq(t, result, i)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for s := range map[string]uint{
			"0": 0,
			"1": 1,
			"2": 2,
			"3": 3,
			"4": 4,
			"5": 5,
			"6": 6,
			"7": 7,
			"8": 8,
			"9": 9,
		} {
			result, err := ParseString(s, UnsignedN[uint](2))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for s, i := range map[string]uint{
			"10":     10,
			"500":    50,
			"100500": 10,
			"10asd":  10,
		} {
			result, err := ParseString(s, UnsignedN[uint](2))
			Check(t, err)
			AssertEq(t, result, i)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for s, i := range map[string]uint{
			"10asdasd":  10,
			"50OIUO":    50,
			"10)(*(0))": 10,
			"10asd":     10,
		} {
			result, err := ParseString(s, UnsignedN[uint](2))
			Check(t, err)
			AssertEq(t, result, i)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, s := range []string{
			"asdasd10",
			" 50OIUO",
			"(10)(*(0))",
			"_10asd",
		} {
			result, err := ParseString(s, UnsignedN[uint](2))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}
