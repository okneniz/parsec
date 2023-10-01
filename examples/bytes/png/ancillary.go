package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

type Ancillary struct {
	length    uint32
	chunkType string
	data      []byte
	crc       uint32
}

func (c *Ancillary) Length() uint32 {
	return c.length
}

func (c *Ancillary) Type() string {
	return c.chunkType
}

func (c *Ancillary) Data() []byte {
	return c.data
}

func (c *Ancillary) CRC() uint32 {
	return c.crc
}

func (c *Ancillary) String() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf("\t length: %v\n", c.length))
	b.WriteString(fmt.Sprintf("\t data: %v\n", c.data))
	b.WriteString(fmt.Sprintf("\t crc: %v\n", c.crc))

	return b.String()
}

func AncillaryChunk(
	size uint32,
	chunkType string,
) p.Combinator[byte, int, *Ancillary] {
	return func(buffer p.Buffer[byte, int]) (*Ancillary, error) {
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

		return &Ancillary{
			length:    size,
			chunkType: chunkType,
			data:      data,
			crc:       crc,
		}, nil
	}
}
