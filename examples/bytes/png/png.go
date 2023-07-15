package png

import (
	. "git.sr.ht/~okneniz/parsec/bytes"
	p "git.sr.ht/~okneniz/parsec/common"
	"golang.org/x/exp/constraints"

	"bytes"
	"encoding/binary"
)

// https://www.w3.org/TR/png/#11IHDR

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
