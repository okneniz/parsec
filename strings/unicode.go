package strings

import (
	"unicode"

	"github.com/okneniz/parsec"
)

// Control parses a control character as defined by unicode.IsControl.
func Control[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsControl)
}

// Digit parses a decimal digit character as defined by unicode.IsDigit.
func Digit[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsDigit)
}

// Graphic parses a graphic character as defined by unicode.IsGraphic.
func Graphic[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsGraphic)
}

// Letter parses a letter character as defined by unicode.IsLetter.
func Letter[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsLetter)
}

// Lower parses a lower case character as defined by unicode.IsLower.
func Lower[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsLower)
}

// Mark parses a mark character as defined by unicode.IsMark.
func Mark[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsMark)
}

// Number parses a number character as defined by unicode.IsNumber.
func Number[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsNumber)
}

// Print parses a printable character as defined by unicode.IsPrint.
func Print[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsPrint)
}

// Punct parses a punctuation character as defined by unicode.IsPunct.
func Punct[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsPunct)
}

// Space parses a whitespace character as defined by unicode.IsSpace.
func Space[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsSpace)
}

// Symbol parses a symbolic character as defined by unicode.IsSymbol.
func Symbol[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsSymbol)
}

// Title parses a title case character as defined by unicode.IsTitle.
func Title[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsTitle)
}

// Upper parses an upper case character as defined by unicode.IsUpper.
func Upper[S any](errMessage string) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, unicode.IsUpper)
}

// RangeTable parses a character from the given unicode.RangeTable,
// see https://pkg.go.dev/unicode#pkg-constants for the available tables.
func RangeTable[S any](
	errMessage string,
	tbl *unicode.RangeTable,
) parsec.Combinator[rune, Position, rune, S] {
	return parsec.Satisfy[rune, Position, S](errMessage, true, func(x rune) bool {
		return unicode.Is(tbl, x)
	})
}
