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

func TestControl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Control("expected control"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Control("expected control"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Control("expected control"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestDigit(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Digit("expected digit"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Digit("expected digit"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Digit("expected digit"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestGraphic(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Graphic("expected graphic"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Graphic("expected graphic"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Graphic("expected graphic"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Graphic("expected graphic"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestLetter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Letter("expected letter"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), Letter("expected letter"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Letter("expected letter"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Letter("expected letter"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Letter("expected letter"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestLower(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Lower("expected lower"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), Lower("expected lower"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Lower("expected lower"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestMark(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range markChars {
			result, err := ParseString(string(c), Mark("expected mark"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), Mark("expected mark"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Mark("expected mark"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Mark("expected mark"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Mark("expected mark"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestNumber(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Number("expected number"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Number("expected number"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), Number("expected number"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Number("expected number"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Number("expected number"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestPrint(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Print("expected print"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Print("expected print"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range graphChars {
			result, err := ParseString(string(c), Print("expected print"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range controllChars {
			result, err := ParseString(string(c), Print("expected print"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 5", func(t *testing.T) {
		for _, c := range spaceChars {
			if c == ' ' {
				continue
			}
			result, err := ParseString(string(c), Print("expected print"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestPunct(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range puntcChars {
			result, err := ParseString(string(c), Punct("expected punctuation"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range spaceChars {
			result, err := ParseString(string(c), Punct("expected punctuation"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Punct("expected punctuation"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 4", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Punct("expected punctuation"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestSpace(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range spaceChars {
			result, err := ParseString(string(c), Space("expected space"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Space("expected space"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Space("expected space"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestSymbol(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range symbolChars {
			result, err := ParseString(string(c), Symbol("expected symbol"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Symbol("expected symbol"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Symbol("expected symbol"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestTitle(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range titleChars {
			result, err := ParseString(string(c), Title("expected title"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Title("expected title"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Title("expected title"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}

func TestUpper(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		for _, c := range upperLettersChars {
			result, err := ParseString(string(c), Upper("expected upper"))
			Check(t, err)
			AssertEq(t, result, c)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		for _, c := range lettersChars {
			result, err := ParseString(string(c), Upper("expected upper"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		for _, c := range digitsChars {
			result, err := ParseString(string(c), Upper("expected upper"))
			AssertError(t, err)
			AssertEq(t, result, 0)
		}
	})
}
