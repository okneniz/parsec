package bytes

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

type buffer struct {
	data     []byte
	position int
}

func (s *buffer) Read(x bool) (byte, error) {
	if s.position >= len(s.data) {
		return 0, p.EndOfFile
	}

	b := s.data[s.position]
	if x {
		s.position++
	}

	return b, nil
}

func (s *buffer) Seek(x int) {
	s.position = x
}

func (s *buffer) Position() int {
	return s.position
}

func (s *buffer) IsEOF() bool {
	return s.position >= len(s.data)
}

func Buffer(data []byte) *buffer {
	b := new(buffer)
	b.data = data
	b.position = 0
	return b
}
