package message_pack

import (
	"encoding/binary"

	b "git.sr.ht/~okneniz/parsec/bytes"
	c "git.sr.ht/~okneniz/parsec/common"
)

// https://github.com/msgpack/msgpack/blob/master/spec.md

func MessagePack() c.Combinator[byte, int, Type] {
	primitive := c.Try(primitiveParser())

	return c.Choice[byte, int, Type](
		primitive,
	)
}

func primitiveParser() c.Combinator[byte, int, Type] {
	cases := map[byte]c.Combinator[byte, int, Type]{
		0xc0: Const[Type](Nil{}),
		0xc2: Const[Type](Boolean(false)),
		0xc3: Const[Type](Boolean(true)),

		0xc4: binaryParser[uint8](1),
		0xc5: binaryParser[uint16](2),
		0xc6: binaryParser[uint32](4),

		0xca: b.Cast(
			b.ReadAs[float32](4, binary.BigEndian),
			func(f float32) (Type, error) {
				return Float32(f), nil
			},
		),
		0xcb: b.Cast(
			b.ReadAs[float64](8, binary.BigEndian),
			func(f float64) (Type, error) {
				return Float64(f), nil
			},
		),

		0xcc: b.Cast(
			b.ReadAs[uint8](1, binary.BigEndian),
			func(f uint8) (Type, error) {
				return Unsigned8(f), nil
			},
		),
		0xcd: b.Cast(
			b.ReadAs[uint16](2, binary.BigEndian),
			func(f uint16) (Type, error) {
				return Unsigned16(f), nil
			},
		),
		0xce: b.Cast(
			b.ReadAs[uint32](4, binary.BigEndian),
			func(f uint32) (Type, error) {
				return Unsigned32(f), nil
			},
		),
		0xcf: b.Cast(
			b.ReadAs[uint64](8, binary.BigEndian),
			func(f uint64) (Type, error) {
				return Unsigned64(f), nil
			},
		),

		0xd0: b.Cast(
			b.ReadAs[int8](1, binary.BigEndian),
			func(f int8) (Type, error) {
				return Signed8(f), nil
			},
		),
		0xd1: b.Cast(
			b.ReadAs[int16](2, binary.BigEndian),
			func(f int16) (Type, error) {
				return Signed16(f), nil
			},
		),
		0xd2: b.Cast(
			b.ReadAs[int32](4, binary.BigEndian),
			func(f int32) (Type, error) {
				return Signed32(f), nil
			},
		),
		0xd3: b.Cast(
			b.ReadAs[int64](8, binary.BigEndian),
			func(f int64) (Type, error) {
				return Signed64(f), nil
			},
		),

		// // TODO : not big-endian?
		0xd9: stringParser[uint8](1), // is a 8-bit unsigned integer which represents N
		0xda: stringParser[uint16](2),
		0xdb: stringParser[uint32](4),
	}

	return MapAs[byte, int, byte, Type](cases, b.Any())
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

func Const[S any](value S) c.Combinator[byte, int, S] {
	return func(_ c.Buffer[byte, int]) (S, error) {
		return value, nil
	}
}

func stringParser[T b.Number](size int) c.Combinator[byte, int, Type] {
	comb := b.ReadAs[T](size, binary.BigEndian)

	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := comb(buffer)
		if err != nil {
			return String(""), err
		}

		data := make([]byte, int(size))
		for i := 0; i < int(size); i++ {
			x, err := buffer.Read(true)
			if err != nil {
				return String(""), err
			}

			data[i] = x
		}

		return String(string(data)), nil
	}
}

func binaryParser[T b.Number](size int) c.Combinator[byte, int, Type] {
	comb := b.ReadAs[T](size, binary.BigEndian)

	return func(buffer c.Buffer[byte, int]) (Type, error) {
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
		return nil, nil
	}
}

func mapParser[K comparable, V any]() c.Combinator[byte, int, Map[K, V]] {
	return func(buffer c.Buffer[byte, int]) (Map[K, V], error) {
		return nil, nil
	}
}
