package png

import (
	"encoding/binary"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
)

// https://www.w3.org/TR/png

func PNG() parsec.Combinator[byte, int, *File, parsec.Stateless] {
	parseHeader := parsec.SkipSequenceOf[byte, int, byte, parsec.Stateless](
		"expected PNG header",
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	)

	parseChunks := bytes.Some(
		1,
		"expected PNG chunks",
		bytes.Try(parseChunk()),
	)

	return func(state parsec.Stateless, buffer parsec.Buffer[byte, int]) (*File, parsec.Error[int]) {
		_, err := parseHeader(state, buffer)
		if err != nil {
			return nil, err
		}

		chunks, err := parseChunks(state, buffer)
		if err != nil {
			return nil, err
		}

		return &File{chunks}, nil
	}
}

func parseChunk() parsec.Combinator[byte, int, Chunk, parsec.Stateless] {
	lenghtOfChunk := bytes.ReadAs[uint32, parsec.Stateless](
		4,
		"expected 4 bytes (uint32) of chunk length",
		binary.BigEndian,
	)

	typeOfChunk := bytes.Count(
		4,
		"expecte 4 bytes of chunk type",
		bytes.Any[parsec.Stateless](),
	)

	return func(state parsec.Stateless, buffer parsec.Buffer[byte, int]) (Chunk, parsec.Error[int]) {
		length, err := lenghtOfChunk(state, buffer)
		if err != nil {
			return nil, err
		}

		chunkType, err := typeOfChunk(state, buffer)
		if err != nil {
			return nil, err
		}

		chunkTypeString := string(chunkType)

		switch chunkTypeString {
		case "IHDR":
			return IHDRChunk(length)(state, buffer)
		case "PLTE":
			return PLTEChunk(length)(state, buffer)
		case "IDAT":
			return IDATChunk(length)(state, buffer)
		case "IEND":
			return IENDChunk(length)(state, buffer)
		default:
			return AncillaryChunk(length, chunkTypeString)(state, buffer)
		}
	}
}
