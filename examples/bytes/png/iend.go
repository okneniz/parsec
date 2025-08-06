package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

type IEND struct {
	length uint32
	crc    uint32
}

func (c *IEND) Length() uint32 {
	return c.length
}

func (c *IEND) Type() string {
	return "IEND"
}

func (c *IEND) Data() []byte {
	return nil
}

func (c *IEND) CRC() uint32 {
	return c.crc
}

func (c *IEND) String() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf("\t length: %v\n", c.length))
	b.WriteString(fmt.Sprintf("\t crc: %v\n", c.crc))

	return b.String()
}

func IENDChunk(size uint32) common.Combinator[byte, int, *IEND] {
	return func(buffer common.Buffer[byte, int]) (*IEND, common.Error[int]) {
		pos := buffer.Position()
		if size > 0 {
			return nil, common.NewParseError(
				pos, fmt.Sprintf("IEND must be empty chunk, but actual %v", size),
			)
		}

		crc, err := ReadAs[uint32](4, "4 bytes of CRC", binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		return &IEND{
			length: size,
			crc:    crc,
		}, nil
	}
}
