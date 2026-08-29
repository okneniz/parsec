package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
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

	fmt.Fprintf(b, "\t length: %v\n", c.length)
	// fmt.Fprintf(b, "\t data: %v\n", c.data)
	fmt.Fprintf(b, "\t crc: %v\n", c.crc)

	return b.String()
}

func IDATChunk(size uint32) parsec.Combinator[byte, int, *IDAT] {
	parseData := bytes.Count[byte](
		int(size),
		fmt.Sprintf("expected %d bytes of IDAT chunk", size),
		bytes.Any(),
	)

	parseCRC := bytes.ReadAs[uint32](4, "expected four bytes of CRC", binary.BigEndian)

	return func(buffer parsec.Buffer[byte, int]) (*IDAT, parsec.Error[int]) {
		var data []byte
		var err parsec.Error[int]

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
