package parsec

type Buffer[T any] interface {
	Read(bool) (T, bool)
	Seek(int)
	Position() int
}

type bytesBuffer struct {
	data []byte
	position int
}

func (s *bytesBuffer) Read(x bool) (byte, bool) {
	if s.position >= len(s.data) {
		return 0, false
	}

	b := s.data[s.position]
	s.position++

	return b, true
}

func (s *bytesBuffer) Seek(x int) {
	s.position = x
}

func (s *bytesBuffer) Position() int {
	return s.position
}

func BytesBuffer(data []byte) *bytesBuffer {
	b := new(bytesBuffer)
	b.data = data
	b.position = 0
	return b
}
