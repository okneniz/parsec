package strings

import (
	"unicode"

	"github.com/okneniz/parsec/common"
)

// Control - parse control UTF-8 characters.
func Control(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsControl)
}

// Digit - parse decimal digit UTF-8 characters.
func Digit(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsDigit)
}

// Graphic - parse graphic UTF-8 characters.
func Graphic(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsGraphic)
}

// Letter - parse letter UTF-8 characters.
func Letter(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLetter)
}

// Lower - parse UTF-8 character in lower case.
func Lower(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLower)
}

// Lower - parse mark UTF-8 characters.
func Mark(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsMark)
}

// Number - parse UTF-8 number characters.
func Number(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsNumber)
}

// Print - parse printable UTF-8 characters.
func Print(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPrint)
}

// Punct - parse UTF-8 punctuation character.
func Punct(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPunct)
}

// Space - parse UTF-8 space character.
func Space(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSpace)
}

// Space - parse UTF-8 symbolic character.
func Symbol(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSymbol)
}

// Title - parse UTF-8 character in title case.
func Title(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsTitle)
}

// Lower - parse UTF-8 character in upper case.
func Upper(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsUpper)
}
