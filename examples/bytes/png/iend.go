package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

type IEND struct {
	length uint32
	data   []byte
	crc    uint32
}

func (c *IEND) Length() uint32 {
	return c.length
}

func (c *IEND) Type() string {
	return "IEND"
}

func (c *IEND) Data() []byte {
	return c.data
}

func (c *IEND) CRC() uint32 {
	return c.crc
}

func (c *IEND) String() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf("\t length: %v\n", c.length))
	b.WriteString(fmt.Sprintf("\t data: %v\n", c.data))
	b.WriteString(fmt.Sprintf("\t crc: %v\n", c.crc))

	return b.String()
}

func IENDChunk(size uint32) p.Combinator[byte, int, *IEND] {
	return func(buffer p.Buffer[byte, int]) (*IEND, error) {
		var data []byte
		var err error

		if size > 0 {
			return nil, fmt.Errorf("IEND must be empty chunk, but actual %v", size)
		}

		crc, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		return &IEND{
			length: size,
			data:   data,
			crc:    crc,
		}, nil
	}
}
