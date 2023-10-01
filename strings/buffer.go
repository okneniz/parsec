package strings

import (
	"github.com/okneniz/parsec/common"
)

var (
	defaultNewLineRunes = map[rune]struct{}{'\n': {}}
)

type buffer struct {
	data         []rune
	position     Position
	newLineRunes map[rune]struct{}
}

// Read - read next item, if greedy buffer keep position after reading.
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

// Seek - change buffer position
func (b *buffer) Seek(x Position) {
	b.position = x
}

// Position - return current buffer position
func (b *buffer) Position() Position {
	return b.position
}

// IsEOF - true if buffer ended.
func (b *buffer) IsEOF() bool {
	return b.position.index >= len(b.data)
}

// Buffer - make buffer which can read text on input and use
// struct for positions.
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
