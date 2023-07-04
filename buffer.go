package parsec

type Buffer[T any, S any] interface {
	Read(bool) (T, error)
	Seek(S)
	Position() S
	IsEOF() bool
}

type bytesBuffer struct {
	data     []byte
	position int
}

func (s *bytesBuffer) Read(x bool) (byte, error) {
	if s.position >= len(s.data) {
		return 0, EndOfFile
	}

	b := s.data[s.position]
	if x {
		s.position++
	}

	return b, nil
}

func (s *bytesBuffer) Seek(x int) {
	s.position = x
}

func (s *bytesBuffer) Position() int {
	return s.position
}

func (s *bytesBuffer) IsEOF() bool {
	return s.position >= len(s.data)
}

func BytesBuffer(data []byte) *bytesBuffer {
	b := new(bytesBuffer)
	b.data = data
	b.position = 0
	return b
}
