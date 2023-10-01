package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

var (
	// source - https://www.utf8-chartable.de/unicode-utf8-table.pl?start=9472&unicodeinhtml=dec

	controllChars = []rune{
		'\u0000', // NULL
		'\u0009', // HORIZONTAL TABULATION
		'\u000A', // LINE FEED
		'\u000C', // FORM FEED
		'\u000D', // CARRIAGE RETURN
		'\u0085', // NEXT LINE
	}

	graphChars = []rune{
		'\u250C', // ┌
		'\u252B', // ┫
		'\u2547', // ╇
		'\u2564', // ╤
		'\u2573', // ╳
		'\u2593', // ▓
	}

	spaceChars        = []rune("\t\n\v\f\r \u0085\u00A0")
	puntcChars        = []rune(`!"#%&()*,-./:;?@[\]_{}¡§«¶·»¿;·`)
	digitsChars       = []rune("0123456789")
	lettersChars      = []rune("abcdefghijklmnopqrstuvwxyz")
	upperLettersChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	titleChars        = []rune("ǅǈǋǲ")
	markChars         = []rune{'̳', '̴', '̵', '̶'}
	symbolChars       = []rune("$+<=>^`|~¢£¤¥¦¨©")
)

func TestIsControl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsControl())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsControl())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsControl())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsDigit(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsDigit())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsDigit())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsDigit())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsGraphic(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsGraphic())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsGraphic())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsGraphic())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsGraphic())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsLetter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsLetter())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), IsLetter())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsLetter())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsLetter())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsLetter())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsLower(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsLower())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), IsLower())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsLower())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsMark(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range markChars {
			result, err := ParseString(string(c), IsMark())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), IsMark())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsMark())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsMark())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsMark())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsNumber(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsNumber())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsNumber())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), IsNumber())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsNumber())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsNumber())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsPrint(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsPrint())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsPrint())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), IsPrint())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), IsPrint())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range spaceChars {
			if c == ' ' {
				continue
			}
			result, err := ParseString(string(c), IsPrint())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsPunct(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range puntcChars {
			result, err := ParseString(string(c), IsPunct())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range spaceChars {
			result, err := ParseString(string(c), IsPunct())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsPunct())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsPunct())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsSpace(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range spaceChars {
			result, err := ParseString(string(c), IsSpace())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsSpace())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsSpace())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsSymbol(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range symbolChars {
			result, err := ParseString(string(c), IsSymbol())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsSymbol())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsSymbol())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsTitle(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range titleChars {
			result, err := ParseString(string(c), IsTitle())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsTitle())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsTitle())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestIsUpper(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), IsUpper())
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), IsUpper())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), IsUpper())
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}
