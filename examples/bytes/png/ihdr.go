package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
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

func IHDRChunk(size uint32) common.Combinator[byte, int, *IHDR] {
	parseData := bytes.Count[byte](
		int(size),
		fmt.Sprintf("expected %d bytes of IHDR data", size),
		bytes.Any(),
	)

	parseWidth := bytes.ReadAs[uint32](4, "expected 4 bytes of IHDR width", binary.BigEndian)
	parseHeight := bytes.ReadAs[uint32](4, "expected 4 bytes of IHDR height", binary.BigEndian)
	parseBitDepth := bytes.Satisfy("expected 1 byte of IHDR bit depth", true, common.Anything[byte])
	parseColorType := bytes.Satisfy("expected 1 byte of color type", true, common.Anything[byte])
	parseCompressionMethod := bytes.Satisfy("expected 1 byte of compression method", true, common.Anything[byte])
	parseFilterMethod := bytes.Satisfy("expected 1 byte of filter method", true, common.Anything[byte])
	parseInterfaceMethod := bytes.Satisfy("expected 1 byte of interface method", true, common.Anything[byte])
	parseCRC := bytes.ReadAs[uint32](4, "expecte 4 bytes of CRC", binary.BigEndian)

	return func(buffer common.Buffer[byte, int]) (*IHDR, common.Error[int]) {
		pos := buffer.Position()

		data, err := parseData(buffer)
		if err != nil {
			return nil, err
		}

		buffer.Seek(pos)

		width, err := parseWidth(buffer)
		if err != nil {
			return nil, err
		}

		height, err := parseHeight(buffer)
		if err != nil {
			return nil, err
		}

		bitDepth, err := parseBitDepth(buffer)
		if err != nil {
			return nil, err
		}

		colorType, err := parseColorType(buffer)
		if err != nil {
			return nil, err
		}

		compressionMethod, err := parseCompressionMethod(buffer)
		if err != nil {
			return nil, err
		}

		filterMethod, err := parseFilterMethod(buffer)
		if err != nil {
			return nil, err
		}

		interfaceMethod, err := parseInterfaceMethod(buffer)
		if err != nil {
			return nil, err
		}

		crc, err := parseCRC(buffer)
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
