package strings

import (
	"unicode"

	p "github.com/okneniz/parsec/common"
)

// TODO : remove Is prefix?

// IsControl - parse control UTF-8 characters.
func IsControl() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsControl)
}

// IsDigit - parse decimal digit UTF-8 characters.
func IsDigit() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsDigit)
}

// IsGraphic - parse graphic UTF-8 characters.
func IsGraphic() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsGraphic)
}

// IsLetter - parse letter UTF-8 characters.
func IsLetter() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsLetter)
}

// IsLower - parse UTF-8 character in lower case.
func IsLower() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsLower)
}

// IsLower - parse mark UTF-8 characters.
func IsMark() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsMark)
}

// IsNumber - parse UTF-8 number characters.
func IsNumber() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsNumber)
}

// IsPrint - parse printable UTF-8 characters.
func IsPrint() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsPrint)
}

// IsPunct - parse UTF-8 punctuation character.
func IsPunct() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsPunct)
}

// IsSpace - parse UTF-8 space character.
func IsSpace() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsSpace)
}

// IsSpace - parse UTF-8 symbolic character.
func IsSymbol() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsSymbol)
}

// IsTitle - parse UTF-8 character in title case.
func IsTitle() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsTitle)
}

// IsLower - parse UTF-8 character in upper case.
func IsUpper() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsUpper)
}
