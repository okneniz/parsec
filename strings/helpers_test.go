package strings

import (
	"testing"
	. "git.sr.ht/~okneniz/parsec/testing"
)

func TestParens(t *testing.T){
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

func TestBraces(t *testing.T){
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

func TestAngles(t *testing.T){
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

func TestSquares(t *testing.T){
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

func TestSemi(t *testing.T){
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

func TestComma(t *testing.T){
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

func TestColon(t *testing.T){
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

func TestDot(t *testing.T){
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
