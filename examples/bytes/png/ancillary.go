package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
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

	fmt.Fprintf(b, "\t length: %v\n", c.length)
	fmt.Fprintf(b, "\t chunk type: %v\n", c.chunkType)
	fmt.Fprintf(b, "\t data: %v\n", c.data)
	fmt.Fprintf(b, "\t CRC: %v\n", c.crc)

	return b.String()
}

func AncillaryChunk(
	size uint32,
	chunkType string,
) common.Combinator[byte, int, *Ancillary] {
	parseData := bytes.Count[byte](
		int(size),
		fmt.Sprintf("expected ancillary chunk (%s)", chunkType),
		bytes.Any(),
	)

	parseCRC := bytes.ReadAs[uint32](
		4,
		"expected 4 big endian bytes of CRC",
		binary.BigEndian,
	)

	return func(buffer common.Buffer[byte, int]) (*Ancillary, common.Error[int]) {
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

		return &Ancillary{
			length:    size,
			chunkType: chunkType,
			data:      data,
			crc:       crc,
		}, nil
	}
}
