package png

import (
	"encoding/binary"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

// https://www.w3.org/TR/png

func PNG() common.Combinator[byte, int, *File] {
	parseHeader := common.SkipSequenceOf[byte, int, byte](
		"expected PNG header",
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	)

	parseChunks := bytes.Some(
		1,
		"expected PNG chunks",
		bytes.Try(parseChunk()),
	)

	return func(buffer common.Buffer[byte, int]) (*File, common.Error[int]) {
		_, err := parseHeader(buffer)
		if err != nil {
			return nil, err
		}

		chunks, err := parseChunks(buffer)
		if err != nil {
			return nil, err
		}

		return &File{chunks}, nil
	}
}

func parseChunk() common.Combinator[byte, int, Chunk] {
	lenghtOfChunk := bytes.ReadAs[uint32](
		4,
		"expected 4 bytes (uint32) of chunk length",
		binary.BigEndian,
	)

	typeOfChunk := bytes.Count(
		4,
		"expecte 4 bytes of chunk type",
		bytes.Any(),
	)

	return func(buffer common.Buffer[byte, int]) (Chunk, common.Error[int]) {
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
