package strings

import (
	"unicode"

	"github.com/okneniz/parsec/common"
)

// Control - parse control UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Control(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsControl)
}

// Digit - parse decimal digit UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Digit(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsDigit)
}

// Graphic - parse graphic UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Graphic(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsGraphic)
}

// Letter - parse letter UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Letter(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLetter)
}

// Lower - parse UTF-8 character in lower case.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Lower(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsLower)
}

// Lower - parse mark UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Mark(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsMark)
}

// Number - parse UTF-8 number characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Number(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsNumber)
}

// Print - parse printable UTF-8 characters.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Print(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPrint)
}

// Punct - parse UTF-8 punctuation character.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Punct(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsPunct)
}

// Space - parse UTF-8 space character.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Space(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSpace)
}

// Space - parse UTF-8 symbolic character.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Symbol(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsSymbol)
}

// Title - parse UTF-8 character in title case.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Title(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsTitle)
}

// Lower - parse UTF-8 character in upper case.
// Read more about utf characters tables - https://pkg.go.dev/unicode#pkg-constants
func Upper(errMessage string) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, unicode.IsUpper)
}

// RangeTable - parse UTF-8 character in upper case.
// Read more about utf characters tables - https://pkg.go.dev/unicode
func RangeTable(
	errMessage string,
	tbl *unicode.RangeTable,
) common.Combinator[rune, Position, rune] {
	return common.Satisfy[rune, Position](errMessage, true, func(x rune) bool {
		return unicode.Is(tbl, x)
	})
}
