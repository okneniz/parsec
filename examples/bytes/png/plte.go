package png

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

type PLTE struct {
	length uint32
	data   []byte
	crc    uint32

	Entries []*RGB
}

func (c *PLTE) Length() uint32 {
	return c.length
}

func (c *PLTE) Type() string {
	return "PLTE"
}

func (c *PLTE) Data() []byte {
	return c.data
}

func (c *PLTE) CRC() uint32 {
	return c.crc
}

func (c *PLTE) String() string {
	b := new(strings.Builder)

	b.WriteString("\t entries:\n")
	b.WriteString("\t [")
	for i, entry := range c.Entries {
		b.WriteString("\t ")
		b.WriteString(entry.String())
		if len(c.Entries) != i+1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]\n")

	return b.String()
}

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (c *RGB) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.Red, c.Green, c.Blue)
}

func PLTEChunk(size uint32) p.Combinator[byte, int, *PLTE] {
	return func(buffer p.Buffer[byte, int]) (*PLTE, error) {
		pos := buffer.Position()

		data, err := Count[byte](int(size), Any())(buffer)
		if err != nil {
			return nil, err
		}

		buffer.Seek(pos)

		entries := make([]*RGB, 0, uint32(size/3))

		for i := uint32(0); i < size/3; i++ {
			red, err := ReadAs[uint8](1, binary.BigEndian)(buffer)
			if err != nil {
				return nil, err
			}

			green, err := ReadAs[uint8](1, binary.BigEndian)(buffer)
			if err != nil {
				return nil, err
			}

			blue, err := ReadAs[uint8](1, binary.BigEndian)(buffer)
			if err != nil {
				return nil, err
			}

			entries = append(entries, &RGB{
				Red:   red,
				Green: green,
				Blue:  blue,
			})
		}

		// TODO : check end pos?

		buffer.Seek(pos + int(size))

		crc, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		return &PLTE{
			length:  size,
			data:    data,
			crc:     crc,
			Entries: entries,
		}, nil
	}
}
