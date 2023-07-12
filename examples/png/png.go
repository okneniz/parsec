package png

import (
	. "git.sr.ht/~okneniz/parsec/bytes"
	p "git.sr.ht/~okneniz/parsec/common"
	"golang.org/x/exp/constraints"

	"bytes"
	"encoding/binary"

	"fmt"
	"strings"
)

// https://www.w3.org/TR/png/#11IHDR

type File struct {
	chunks []Chunk
}

func (f *File) String() string {
	b := new(strings.Builder)

	for _, chunk := range f.chunks {

		b.WriteString(chunk.Type())
		b.WriteString(":\n")
		b.WriteString(chunk.String())
	}

	return b.String()
}

func PNG() p.Combinator[byte, int, *File] {
	head := SkipSequenceOf[byte, int, byte](0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A)

	return func(buffer p.Buffer[byte, int]) (*File, error) {
		_, err := head(buffer)
		if err != nil {
			return nil, err
		}

		chunks, err := Some(1, Try(chunk()))(buffer)
		if err != nil {
			return nil, err
		}

		return &File{chunks}, nil
	}
}

type Chunk interface {
	Length() uint32
	Type() string
	Data() []byte
	CRC() uint32
	String() string
}

func chunk() p.Combinator[byte, int, Chunk] {
	lenghtOfChunk := ReadAs[uint32](4, binary.BigEndian)
	typeOfChunk := Count(4, Any())

	return func(buffer p.Buffer[byte, int]) (Chunk, error) {
		length, err := lenghtOfChunk(buffer)
		if err != nil {
			return nil, err
		}

		chunkType, err := typeOfChunk(buffer)
		if err != nil {
			return nil, err
		}

		chunkTypeString := string(chunkType)

		switch chunkTypeString {
		case "IHDR":
			return IHDRChunk(length)(buffer)
		case "PLTE":
			return PLTEChunk(length)(buffer)
		case "IDAT":
			return IDATChunk(length)(buffer)
		case "IEND":
			return IENDChunk(length)(buffer)
		default:
			return AncillaryChunk(length, chunkTypeString)(buffer)
		}
	}
}

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
			FilterMethod: 	   filterMethod,
			InterfaceMethod:   interfaceMethod,
		}

		return result, nil
	}
}

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
		if len(c.Entries) != i + 1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]\n")

	return b.String()
}


type RGB struct {
	Red	  uint8
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

		entries := make([]*RGB, 0, uint32(size / 3))

		for i := uint32(0); i < size / 3; i++ {
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
				Red: red,
				Green: green,
				Blue: blue,
			})
		}

		// TODO : check end pos?

		buffer.Seek(pos + int(size))

		crc, err := ReadAs[uint32](4, binary.BigEndian)(buffer)
		if err != nil {
			return nil, err
		}

		return &PLTE{
			length: size,
			data: data,
			crc: crc,
			Entries: entries,
		}, nil
	}
}

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
			length: size,
			chunkType: chunkType,
			data: data,
			crc: crc,
		}, nil
	}
}

type IDAT struct {
	length    uint32
	data      []byte
	crc       uint32
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
			data: data,
			crc: crc,
		}, nil
	}
}

type IEND struct {
	length    uint32
	data      []byte
	crc       uint32
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
			data: data,
			crc: crc,
		}, nil
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

func ReadAs[T Number](
	size int,
	order binary.ByteOrder,
) p.Combinator[byte, int, T] {
	return func(buffer p.Buffer[byte, int]) (T, error) {
		var result T

		input, err := Count(size, Any())(buffer)
		if err != nil {
			return result, err
		}

		buf := bytes.NewReader(input)

		err = binary.Read(buf, order, &result)
		if err != nil {
			return result, err
		}

		return result, nil
	}
}

func SkipSequence[T, P, S any](combs ...p.Combinator[T, P, S]) p.Combinator[T, P, S] {
	return func(buffer p.Buffer[T, P]) (S, error) {
		var result S

		for _, c := range combs {
			_, err := c(buffer)
			if err != nil {
				return result, err
			}
		}

		return result, nil
	}
}

func SkipSequenceOf[T comparable, P, S any](data ...T) p.Combinator[T, P, S] {
	return func(buffer p.Buffer[T, P]) (S, error) {
		var result S

		for _, x := range data {
			r, err := buffer.Read(true)
			if err != nil {
				return result, err
			}
			if x != r {
				return result, p.NothingMatched
			}
		}

		return result, nil
	}
}
