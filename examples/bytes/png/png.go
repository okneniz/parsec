package png

import (
	"encoding/binary"
	. "github.com/okneniz/parsec/bytes"
	p "github.com/okneniz/parsec/common"
)

// https://www.w3.org/TR/png

func PNG() p.Combinator[byte, int, *File] {
	head := p.SkipSequenceOf[byte, int, byte](0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A)

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
