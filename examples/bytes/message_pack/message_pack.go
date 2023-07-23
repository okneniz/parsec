package message_pack

import (
	"encoding/binary"

	b "git.sr.ht/~okneniz/parsec/bytes"
	c "git.sr.ht/~okneniz/parsec/common"
	// "encoding/binary"
)

// https://github.com/msgpack/msgpack/blob/master/spec.md

func MessagePack() c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		return nil, nil
	}
}

func primitiveParser() c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		cases := map[byte]c.Combinator[byte, int, Type]{
			0xc0: Const(Nil{}),
			0xc2: Const(Boolean(false)),
			0xc3: Const(Boolean(true)),

			0xc4: binaryParser[uint8](1),
			0xc5: binaryParser[uint16](2),
			0xc6: binaryParser[uint32](4),

			0xca: b.ReadAs[Float32](4, binary.BigEndian),
			0xcb: b.ReadAs[Float64](16, binary.BigEndian),

			0xcc: b.ReadAs[Unsigned8](1, binary.BigEndian),
			0xcd: b.ReadAs[Unsigned16](2, binary.BigEndian),
			0xce: b.ReadAs[Unsigned32](4, binary.BigEndian),
			0xcf: b.ReadAs[Unsigned64](16, binary.BigEndian),

			0xd0: b.ReadAs[Signed8](1, binary.BigEndian),
			0xd1: b.ReadAs[Signed16](2, binary.BigEndian),
			0xd2: b.ReadAs[Signed32](4, binary.BigEndian),
			0xd3: b.ReadAs[Signed64](16, binary.BigEndian),

			// TODO : not big-endian?
			0xd9: stringParser[uint8](1), // is a 8-bit unsigned integer which represents N
			0xda: stringParser[uint16](2),
			0xdb: stringParser[uint32](4),
		}

		return MapAs[byte, int, byte, Type](cases, b.Any())
	}
}

func MapAs[T any, P any, K comparable, V any](
	cases map[K]c.Combinator[T, P, V],
	comb c.Combinator[T, P, K],
) c.Combinator[T, P, V] {
	return func(buffer c.Buffer[T, P]) (V, error) {
		var v V

		key, err := comb(buffer)
		if err != nil {
			return v, err
		}

		parseValue, exists := cases[key]
		if !exists {
			return v, c.NothingMatched
		}

		return parseValue(buffer)
	}
}

func Const[T, P, S any](value S) c.Combinator[T, P, S] {
	return func(_ c.Buffer[T, P]) (S, error) {
		return value, nil
	}
}

func stringParser[T b.Number](size int) c.Combinator[byte, int, String] {
	comb := b.ReadAs[T](size, binary.BigEndian)

	return func(buffer c.Buffer[byte, int]) (String, error) {
		size, err := comb(buffer)
		if err != nil {
			return "", err
		}

		data := make([]byte, int(size))
		for i := 0; i < int(size); i++ {
			x, err := buffer.Read(true)
			if err != nil {
				return "", err
			}

			data[i] = x
		}

		return String(data), nil
	}
}

func binaryParser[T b.Number](size int) c.Combinator[byte, int, Binary] {
	comb := b.ReadAs[T](size, binary.BigEndian)

	return func(buffer c.Buffer[byte, int]) (Binary, error) {
		size, err := comb(buffer)
		if err != nil {
			return nil, err
		}

		data := make([]byte, int(size))
		for i := 0; i < int(size); i++ {
			x, err := buffer.Read(true)
			if err != nil {
				return nil, err
			}

			data[i] = x
		}

		return Binary(data), nil
	}
}


func arrayParser[T any]() c.Combinator[byte, int, Array[T]] {
	return func(buffer c.Buffer[byte, int]) (Array[T], error) {
		return []T{}, nil
	}
}

func mapParser[K comparable, V any]() c.Combinator[byte, int, Map[K, V]] {
	return func(buffer c.Buffer[byte, int]) (Map[K, V], error) {
		return map[K]V{}, nil
	}
}
