package strings

import (
	"unicode"

	"github.com/okneniz/parsec/common"
)

// Control parses a control character as defined by unicode.IsControl.
func Control(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsControl)
}

// Digit parses a decimal digit character as defined by unicode.IsDigit.
func Digit(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsDigit)
}

// Graphic parses a graphic character as defined by unicode.IsGraphic.
func Graphic(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsGraphic)
}

// Letter parses a letter character as defined by unicode.IsLetter.
func Letter(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLetter)
}

// Lower parses a lower case character as defined by unicode.IsLower.
func Lower(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLower)
}

// Mark parses a mark character as defined by unicode.IsMark.
func Mark(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsMark)
}

// Number parses a number character as defined by unicode.IsNumber.
func Number(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsNumber)
}

// Print parses a printable character as defined by unicode.IsPrint.
func Print(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPrint)
}

// Punct parses a punctuation character as defined by unicode.IsPunct.
func Punct(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPunct)
}

// Space parses a whitespace character as defined by unicode.IsSpace.
func Space(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSpace)
}

// Symbol parses a symbolic character as defined by unicode.IsSymbol.
func Symbol(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSymbol)
}

// Title parses a title case character as defined by unicode.IsTitle.
func Title(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsTitle)
}

// Upper parses an upper case character as defined by unicode.IsUpper.
func Upper(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsUpper)
}

// RangeTable parses a character from the given unicode.RangeTable,
// see https://pkg.go.dev/unicode#pkg-constants for the available tables.
func RangeTable(
	errMessage string,
	tbl *unicode.RangeTable,
) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, func(x rune) bool {
		return unicode.Is(tbl, x)
	})
}
