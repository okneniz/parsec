package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
	"unicode"
)

func IsControl() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsControl)
}

func IsDigit() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsDigit)
}

func IsGraphic() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsGraphic)
}

func IsLetter() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsLetter)
}

func IsLower() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsLower)
}

func IsMark() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsMark)
}

func IsNumber() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsNumber)
}

func IsPrint() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsPrint)
}

func IsPunct() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsPunct)
}

func IsSpace() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsSpace)
}

func IsSymbol() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsSymbol)
}

func IsTitle() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsTitle)
}

func IsUpper() p.Combinator[rune, Position, rune] {
	return p.Satisfy[rune, Position](true, unicode.IsUpper)
}
