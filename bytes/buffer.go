package bytes

import (
	"os"

	"github.com/okneniz/parsec"
)

type buffer struct {
	data     []byte
	position int
}

var _ parsec.Buffer[byte, int] = new(buffer)

// Read returns the next byte.
// If greedy is true, the buffer advances and keeps the new position,
// including the case when the calling combinator later fails on this byte;
// with greedy false the byte is only peeked and the position stays unchanged.
func (s *buffer) Read(greedy bool) (byte, error) {
	if s.position >= len(s.data) {
		return 0, parsec.ErrEndOfFile
	}

	b := s.data[s.position]
	if greedy {
		s.position++
	}

	return b, nil
}

// Seek moves the buffer to a previously obtained position.
// It returns parsec.ErrOutOfBounds when the position is invalid.
// Seeking to the current position does nothing.
func (s *buffer) Seek(x int) error {
	if s.position == x {
		return nil
	}

	if x < 0 {
		return parsec.ErrOutOfBounds
	}

	if x >= len(s.data) {
		return parsec.ErrOutOfBounds
	}

	s.position = x
	return nil
}

// Position returns the current buffer position.
func (s *buffer) Position() int {
	return s.position
}

// IsEOF reports whether the buffer is fully consumed.
func (s *buffer) IsEOF() bool {
	return s.position >= len(s.data)
}

// Buffer creates a buffer which reads data and uses
// an integer byte offset as the position.
func Buffer(data []byte) *buffer {
	b := new(buffer)
	b.data = data
	b.position = 0
	return b
}

// BufferFromFile reads the file and creates a buffer over its content.
func BufferFromFile(path string) (*buffer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Buffer(data), nil
}
