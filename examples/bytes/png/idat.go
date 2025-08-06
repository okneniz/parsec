package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
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
	// b.WriteString(fmt.Sprintf("\t data: %v\n", c.data))
	b.WriteString(fmt.Sprintf("\t crc: %v\n", c.crc))

	return b.String()
}

func IDATChunk(size uint32) common.Combinator[byte, int, *IDAT] {
	parseData := bytes.Count[byte](
		int(size),
		fmt.Sprintf("expected %d bytes of IDAT chunk", size),
		bytes.Any(),
	)

	parseCRC := bytes.ReadAs[uint32](4, "expected four bytes of CRC", binary.BigEndian)

	return func(buffer common.Buffer[byte, int]) (*IDAT, common.Error[int]) {
		var data []byte
		var err common.Error[int]

		if size > 0 {
			data, err = parseData(buffer)
			if err != nil {
				return nil, err
			}
		}

		crc, err := parseCRC(buffer)
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
