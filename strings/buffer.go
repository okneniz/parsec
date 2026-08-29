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

// Read returns the next rune.
// If greedy is true, the buffer advances and keeps the new position,
// including the case when the calling combinator later fails on this rune;
// with greedy false the rune is only peeked and the position stays unchanged.
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

// Seek moves the buffer to a previously obtained position.
// It returns common.ErrOutOfBounds when the position is invalid.
// Seeking to the current position does nothing.
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

// Position returns the current buffer position.
func (b *buffer) Position() Position {
	return b.position
}

// IsEOF reports whether the buffer is fully consumed.
func (b *buffer) IsEOF() bool {
	return b.position.index >= len(b.data)
}

// Buffer creates a buffer which reads data and tracks positions
// as line/column/index. The newLineRunes arguments specify which
// characters are treated as line breaks; by default it is '\n'.
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
