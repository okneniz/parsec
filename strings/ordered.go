package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Range(from rune, to rune) p.Combinator[rune, Position, rune] {
	return p.Range[rune, Position](from, to)
}

func NotRange(from rune, to rune) p.Combinator[rune, Position, rune] {
	return p.NotRange[rune, Position](from, to)
}
