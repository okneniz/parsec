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

var _ common.Buffer[rune, Position] = new(buffer)

// Read - read next item, if greedy buffer keep position after reading.
func (b *buffer) Read(greedy bool) (rune, error) {
	if b.IsEOF() {
		return 0, common.ErrEndOfFile
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
// change nothing if you try to seek to the same position
func (b *buffer) Seek(x Position) error {
	if b.position.index == x.index {
		return nil
	}

	if x.index < 0 {
		return common.ErrOutOfBounds
	}

	if x.index >= len(b.data) {
		return common.ErrOutOfBounds
	}

	b.position = x
	return nil
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
