package message_pack

import (
	"encoding/binary"
	"errors"

	b "github.com/okneniz/parsec/bytes"
	c "github.com/okneniz/parsec/common"
)

var (
	ImposibleDataType = errors.New("0xc1 - impossible data type")
)

// https://github.com/msgpack/msgpack/blob/master/spec.md
func MessagePack() c.Combinator[byte, int, Type] {
	cases := map[byte]c.Combinator[byte, int, Type]{
		0xc0: b.Const[Type](Nil{}),

		0xc1: b.Fail[Type](ImposibleDataType),

		// bool
		0xc2: b.Const[Type](Boolean(false)),
		0xc3: b.Const[Type](Boolean(true)),

		// bin
		0xc4: binaryParser[uint8](1),
		0xc5: binaryParser[uint16](2),
		0xc6: binaryParser[uint32](4),

		// ext
		0xc7: extParser(b.ReadAs[int8](1, binary.BigEndian)),
		0xc8: extParser(b.ReadAs[int16](2, binary.BigEndian)),
		0xc9: extParser(b.ReadAs[int32](3, binary.BigEndian)),

		// float
		0xca: b.Cast(
			b.ReadAs[Float32](4, binary.BigEndian),
			func(f Float32) (Type, error) {
				return f, nil
			},
		),
		0xcb: b.Cast(
			b.ReadAs[Float64](8, binary.BigEndian),
			func(f Float64) (Type, error) {
				return f, nil
			},
		),

		// uint
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

		// int
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

		// fixext
		0xd4: extParser(b.Const[uint8](1)),
		0xd5: extParser(b.Const[uint8](2)),
		0xd6: extParser(b.Const[uint8](4)),
		0xd7: extParser(b.Const[uint8](8)),
		0xd8: extParser(b.Const[uint8](16)),

		// strings
		0xd9: stringParser(b.ReadAs[uint8](1, binary.BigEndian)),
		0xda: stringParser(b.ReadAs[uint16](2, binary.BigEndian)),
		0xdb: stringParser(b.ReadAs[uint32](4, binary.BigEndian)),
	}

	// positive fixint
	for i := byte(0x00); i <= byte(0x7f); i++ {
		cases[i] = b.Const[Type](Unsigned8(i))
	}

	// fixstring parser
	for i := byte(0xa0); i <= byte(0xbf); i++ {
		cases[i] = stringParser(
			b.Const[int](int(i - 0xa0)),
		)
	}

	// negative fixint
	for i, z := byte(0xe0), -32; i <= byte(0xff); i, z = i+1, z+1 {
		cases[i] = b.Const[Type](Signed8(z))
		if i == 0xff { // avoid endless loop
			break
		}
	}

	valuesParser := c.MapAs[byte, int, byte, Type](cases, b.Any())

	// fixarray parser
	for i := byte(0x90); i <= byte(0x9f); i++ {
		cases[i] = arrayParser(
			b.Const[int](int(i-0x90)),
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
			b.Const[int](int(i-0x80)),
			valuesParser,
		)
	}

	// map 16 parser
	cases[0xde] = mapParser(
		b.ReadAs[uint16](2, binary.BigEndian),
		valuesParser,
	)

	// map 32 parser
	cases[0xdf] = mapParser(
		b.ReadAs[uint32](4, binary.BigEndian),
		valuesParser,
	)

	return valuesParser
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

func mapParser[T b.Number](
	parseSize c.Combinator[byte, int, T],
	parseValue c.Combinator[byte, int, Type],
) c.Combinator[byte, int, Type] {
	return func(buffer c.Buffer[byte, int]) (Type, error) {
		size, err := parseSize(buffer)
		if err != nil {
			return nil, err
		}

		data := make(Map, int(size))

		for i := 0; i < int(size); i++ {
			key, err := parseValue(buffer)
			if err != nil {
				return nil, err
			}

			value, err := parseValue(buffer)
			if err != nil {
				return nil, err
			}

			data[i] = Pair{
				Key:   key,
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
