package png

import (
	"strings"
)

type File struct {
	chunks []Chunk
}

type Chunk interface {
	Length() uint32
	Type() string
	Data() []byte
	CRC() uint32
	String() string
}

func (f *File) String() string {
	b := new(strings.Builder)

	for _, chunk := range f.chunks {

		b.WriteString(chunk.Type())
		b.WriteString(":\n")
		b.WriteString(chunk.String())
	}

	return b.String()
}
