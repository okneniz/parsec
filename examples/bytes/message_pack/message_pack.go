package message_pack

import (
	"encoding/binary"
	"errors"
	// "fmt"

	b "git.sr.ht/~okneniz/parsec/bytes"
	c "git.sr.ht/~okneniz/parsec/common"
)

// https://github.com/msgpack/msgpack/blob/master/spec.md

func MessagePack() c.Combinator[byte, int, Type] {
	cases := map[byte]c.Combinator[byte, int, Type]{
		0xc2: Const[Type](Boolean(false)),
		0xc3: Const[Type](Boolean(true)),

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

		0xd9: stringParser(
			b.ReadAs[uint8](1, binary.BigEndian),
		),
		0xda: stringParser(
			b.ReadAs[uint16](2, binary.BigEndian),
		),
		0xdb: stringParser(
			b.ReadAs[uint32](4, binary.BigEndian),
		),
	}

	// fixstring parser
	for i := byte(0xa0); i <= byte(0xbf); i++ {
		cases[i] = stringParser(
			Const[int](int(i - 0xa0)),
		)
	}

	// positive fixint
	for i := byte(0x00); i <= byte(0x7f); i++ {
		cases[i] = Const[Type](Signed8(i))
	}

	// negative fixint
	for i := byte(0xe0); i <= byte(0xff); i++ {
		cases[i] = Const[Type](Signed8(0xe0 - 0xff))
		if i == 0xff { // avoid endless loop
			break
		}
	}

	// complex types

	valuesParser := MapAs[byte, int, byte, Type](cases, b.Any())

	cases[0xc0] = Const[Type](Nil{})

	cases[0xc1] = func(buffer c.Buffer[byte, int]) (Type, error) {
		return nil, errors.New("0xc1 - impossible data type")
	}

	cases[0xc4] = binaryParser[uint8](1)
	cases[0xc5] = binaryParser[uint16](2)
	cases[0xc6] = binaryParser[uint32](4)

	// container types

	// fixarray parser
	for i := byte(0x90); i <= byte(0x9f); i++ {
		cases[i] = arrayParser(
			Const[int](int(i - 0x90)),
			valuesParser,
		)
	}

	// array 16
	cases[0xdc] = arrayParser(
		b.ReadAs[uint16](2, binary.BigEndian),
		valuesParser,
	)

	// array 32
	cases[0xdd] = arrayParser(
		b.ReadAs[uint32](4, binary.BigEndian),
		valuesParser,
	)

	// fixmap parser
	for i := byte(0x80); i <= byte(0x8f); i++ {
		cases[i] = mapParser(
			Const[int](int(i - 0x80)),
			valuesParser,
			valuesParser,
		)
	}

	// map 16 parser
	cases[0xde] = mapParser(
		b.ReadAs[uint16](2, binary.BigEndian),
		valuesParser,
		valuesParser,
	)

	// map 32 parser
	cases[0xdf] = mapParser(
		b.ReadAs[uint32](4, binary.BigEndian),
		valuesParser,
		valuesParser,
	)

	// fixext 1
	cases[0xd4] = extParser(
		Const[int](1),
	)

	// fixext 2
	cases[0xd5] = extParser(
		Const[int](2),
	)

	// fixext 4
	cases[0xd6] = extParser(
		Const[int](4),
	)

	// fixext 8
	cases[0xd7] = extParser(
		Const[int](8),
	)

	// fixext 16
	cases[0xd8] = extParser(
		Const[int](16),
	)

	// ext 8
	cases[0xc7] = extParser(
		b.ReadAs[int8](1, binary.BigEndian),
	)

	// ext 16
	cases[0xc8] = extParser(
		b.ReadAs[int16](2, binary.BigEndian),
	)

	// ext 32
	cases[0xc9] = extParser(
		b.ReadAs[int32](3, binary.BigEndian),
	)

	return valuesParser
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

func stringParser[T b.Number](
	parseSize c.Combinator[byte, int, T],
) c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := parseSize(buffer)
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

func arrayParser[T b.Number](
	parseSize c.Combinator[byte, int, T],
	parseValue c.Combinator[byte, int, Type],
) c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := parseSize(buffer)
		if err != nil {
			return nil, err
		}

		data := make(Array, int(size))
		for i := 0; i < int(size); i++ {
			x, err := parseValue(buffer)
			if err != nil {
				return nil, err
			}

			data[i] = x
		}

		return data, nil
	}
}

func mapParser[K Type, T b.Number](
	parseSize c.Combinator[byte, int, T],
	parseKey c.Combinator[byte, int, K],
	parseValue c.Combinator[byte, int, Type],
) c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := parseSize(buffer)
		if err != nil {
			return nil, err
		}

		data := make(Map, int(size))

		for i := 0; i < int(size); i++ {
			key, err := parseKey(buffer)
			if err != nil {
				return nil, err
			}

			value, err := parseValue(buffer)
			if err != nil {
				return nil, err
			}

			data[i] = Pair {
				Key: key,
				Value: value,
			}
		}

		return data, nil
	}
}

func extParser[T b.Number](
	parseSize c.Combinator[byte, int, T],
) c.Combinator[byte, int, Type] {
	nameParser := b.ReadAs[int8](1, binary.BigEndian)

	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := parseSize(buffer)
		if err != nil {
			return nil, err
		}

		name, err := nameParser(buffer)
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

		return Ext{Name: name, Data: data}, nil
	}
}