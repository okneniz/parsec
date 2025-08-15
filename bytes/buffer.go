package bytes

import (
	"os"

	"github.com/okneniz/parsec/common"
)

type buffer struct {
	data     []byte
	position int
}

var _ common.Buffer[byte, int] = new(buffer)

// Read - read next item, if greedy buffer keep position after reading.
func (s *buffer) Read(greedy bool) (byte, error) {
	if s.position >= len(s.data) {
		return 0, common.ErrEndOfFile
	}

	b := s.data[s.position]
	if greedy {
		s.position++
	}

	return b, nil
}

// Seek - change buffer position
// change nothing if you try to seek to the same position
func (s *buffer) Seek(x int) error {
	if s.position == x {
		return nil
	}

	if x < 0 {
		return common.ErrOutOfBounds
	}

	if x >= len(s.data) {
		return common.ErrOutOfBounds
	}

	s.position = x
	return nil
}

// Position - return current buffer position
func (s *buffer) Position() int {
	return s.position
}

// IsEOF - true if buffer ended.
func (s *buffer) IsEOF() bool {
	return s.position >= len(s.data)
}

// Buffer - make buffer which can read bytes on input and use
// integer for positions.
func Buffer(data []byte) *buffer {
	b := new(buffer)
	b.data = data
	b.position = 0
	return b
}

// Buffer - read file and make buffer which can read bytes on input and use
// integer for positions.
func BufferFromFile(path string) (*buffer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Buffer(data), nil
}
