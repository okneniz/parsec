package strings

import (
	"git.sr.ht/~okneniz/parsec/common"
)

var (
	defaultNewLineRunes = map[rune]struct{}{'\n': {}}
)

type buffer struct {
	data         []rune
	position     Position
	newLineRunes map[rune]struct{}
}

func (b *buffer) Read(greedy bool) (rune, error) {
	if b.IsEOF() {
		return 0, common.EndOfFile
	}

	x := b.data[b.position.index]

	if greedy {
		b.position.index++

		if _, isNewLine := b.newLineRunes[x]; isNewLine {
			b.position.column = 0
			b.position.line++
		} else {
			b.position.column++
		}
	}

	return x, nil
}

func (b *buffer) Seek(x Position) {
	b.position = x
}

func (b *buffer) Position() Position {
	return b.position
}

func (b *buffer) IsEOF() bool {
	return b.position.index >= len(b.data)
}

func Buffer(data []rune, newLineRunes ...rune) *buffer {
	b := new(buffer)
	b.data = data
	b.position = Position{0, 0, 0}

	if len(newLineRunes) == 0 {
		b.newLineRunes = defaultNewLineRunes
	} else {
		b.newLineRunes = make(map[rune]struct{})

		for _, x := range newLineRunes {
			b.newLineRunes[x] = struct{}{}
		}
	}

	return b
}
