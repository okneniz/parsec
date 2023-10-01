package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

type IDAT struct {
	length uint32
	data   []byte
	crc    uint32
}

func (c *IDAT) Length() uint32 {
	return c.length
}

func (c *IDAT) Type() string {
	return "IDAT"
}

func (c *IDAT) Data() []byte {
	return c.data
}

func (c *IDAT) CRC() uint32 {
	return c.crc
}

func (c *IDAT) String() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf("\t length: %v\n", c.length))
	b.WriteString(fmt.Sprintf("\t data: %v\n", c.data))
	b.WriteString(fmt.Sprintf("\t crc: %v\n", c.crc))

	return b.String()
}

func IDATChunk(size uint32) p.Combinator[byte, int, *IDAT] {
	return func(buffer p.Buffer[byte, int]) (*IDAT, error) {
		var data []byte
		var err error

		if size > 0 {
			data, err = Count[byte](int(size), Any())(buffer)
			if err != nil {
				return nil, err
			}
		}

		crc, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		return &IDAT{
			length: size,
			data:   data,
			crc:    crc,
		}, nil
	}
}
