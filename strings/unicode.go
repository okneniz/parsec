package strings

import (
	"unicode"

	"github.com/okneniz/parsec"
)

// Control parses a control character as defined by unicode.IsControl.
func Control(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsControl)
}

// Digit parses a decimal digit character as defined by unicode.IsDigit.
func Digit(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsDigit)
}

// Graphic parses a graphic character as defined by unicode.IsGraphic.
func Graphic(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsGraphic)
}

// Letter parses a letter character as defined by unicode.IsLetter.
func Letter(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsLetter)
}

// Lower parses a lower case character as defined by unicode.IsLower.
func Lower(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsLower)
}

// Mark parses a mark character as defined by unicode.IsMark.
func Mark(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsMark)
}

// Number parses a number character as defined by unicode.IsNumber.
func Number(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsNumber)
}

// Print parses a printable character as defined by unicode.IsPrint.
func Print(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsPrint)
}

// Punct parses a punctuation character as defined by unicode.IsPunct.
func Punct(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsPunct)
}

// Space parses a whitespace character as defined by unicode.IsSpace.
func Space(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsSpace)
}

// Symbol parses a symbolic character as defined by unicode.IsSymbol.
func Symbol(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsSymbol)
}

// Title parses a title case character as defined by unicode.IsTitle.
func Title(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsTitle)
}

// Upper parses an upper case character as defined by unicode.IsUpper.
func Upper(errMessage string) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, unicode.IsUpper)
}

// RangeTable parses a character from the given unicode.RangeTable,
// see https://pkg.go.dev/unicode#pkg-constants for the available tables.
func RangeTable(
	errMessage string,
	tbl *unicode.RangeTable,
) parsec.Combinator[rune, Position, rune] {
	return parsec.Satisfy[rune, Position](errMessage, true, func(x rune) bool {
		return unicode.Is(tbl, x)
	})
}
