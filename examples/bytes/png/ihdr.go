package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

type IHDR struct {
	length uint32
	data   []byte
	crc    uint32

	Width             uint32
	Height            uint32
	BitDepth          byte
	ColorType         byte
	CompressionMethod byte
	FilterMethod      byte
	InterfaceMethod   byte
}

func (c *IHDR) Length() uint32 {
	return c.length
}

func (c *IHDR) Type() string {
	return "IHDR"
}

func (c *IHDR) Data() []byte {
	return c.data
}

func (c *IHDR) CRC() uint32 {
	return c.crc
}

func (c *IHDR) String() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf("\t width: %d\n", c.Width))
	b.WriteString(fmt.Sprintf("\t height: %d\n", c.Height))
	b.WriteString(fmt.Sprintf("\t bit depth: %v\n", c.BitDepth))
	b.WriteString(fmt.Sprintf("\t color type: %v\n", c.ColorType))
	b.WriteString(fmt.Sprintf("\t compression type: %v\n", c.CompressionMethod))
	b.WriteString(fmt.Sprintf("\t interface method: %v\n", c.InterfaceMethod))

	return b.String()
}

func IHDRChunk(size uint32) p.Combinator[byte, int, *IHDR] {
	return func(buffer p.Buffer[byte, int]) (*IHDR, error) {
		pos := buffer.Position()

		data, err := Count[byte](int(size), Any())(buffer)
		if err != nil {
			return nil, err
		}

		buffer.Seek(pos)

		width, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		height, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		bitDepth, err := buffer.Read(true)
		if err != nil {
			return nil, err
		}

		colorType, err := buffer.Read(true)
		if err != nil {
			return nil, err
		}

		compressionMethod, err := buffer.Read(true)
		if err != nil {
			return nil, err
		}

		filterMethod, err := buffer.Read(true)
		if err != nil {
			return nil, err
		}

		interfaceMethod, err := buffer.Read(true)
		if err != nil {
			return nil, err
		}

		crc, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		result := &IHDR{
			length:            uint32(size),
			data:              data,
			crc:               crc,
			Width:             width,
			Height:            height,
			BitDepth:          bitDepth,
			ColorType:         colorType,
			CompressionMethod: compressionMethod,
			FilterMethod:      filterMethod,
			InterfaceMethod:   interfaceMethod,
		}

		return result, nil
	}
}
