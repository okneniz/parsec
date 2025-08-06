package message_pack

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

var (
	ImposibleDataType = errors.New("0xc1 - impossible data type")
)

// https://github.com/msgpack/msgpack/blob/master/spec.md
func MessagePack() common.Combinator[byte, int, Type] {
	cases := map[byte]common.Combinator[byte, int, Type]{
		0xc0: bytes.Const[Type](Nil{}),

		0xc1: bytes.Fail[Type]("impossible data type (0xc1)"),

		// bool
		0xc2: bytes.Const[Type](Boolean(false)),
		0xc3: bytes.Const[Type](Boolean(true)),

		// bin
		0xc4: binaryParser[uint8](1),
		0xc5: binaryParser[uint16](2),
		0xc6: binaryParser[uint32](4),

		// ext
		0xc7: extParser("ext8", bytes.ReadAs[int8](1, "ext8", binary.BigEndian)),
		0xc8: extParser("ext16", bytes.ReadAs[int16](2, "ext16", binary.BigEndian)),
		0xc9: extParser("ext32", bytes.ReadAs[int32](3, "ext32", binary.BigEndian)),

		// float
		0xca: bytes.Cast(
			bytes.ReadAs[Float32](4, "float32", binary.BigEndian),
			func(f Float32) (Type, error) {
				return f, nil
			},
		),
		0xcb: bytes.Cast(
			bytes.ReadAs[Float64](8, "float64", binary.BigEndian),
			func(f Float64) (Type, error) {
				return f, nil
			},
		),

		// uint
		0xcc: bytes.Cast(
			bytes.ReadAs[uint8](1, "uint8", binary.BigEndian),
			func(f uint8) (Type, error) {
				return Unsigned8(f), nil
			},
		),
		0xcd: bytes.Cast(
			bytes.ReadAs[uint16](2, "uint16", binary.BigEndian),
			func(f uint16) (Type, error) {
				return Unsigned16(f), nil
			},
		),
		0xce: bytes.Cast(
			bytes.ReadAs[uint32](4, "uint32", binary.BigEndian),
			func(f uint32) (Type, error) {
				return Unsigned32(f), nil
			},
		),
		0xcf: bytes.Cast(
			bytes.ReadAs[uint64](8, "uint64", binary.BigEndian),
			func(f uint64) (Type, error) {
				return Unsigned64(f), nil
			},
		),

		// int
		0xd0: bytes.Cast(
			bytes.ReadAs[int8](1, "int8", binary.BigEndian),
			func(f int8) (Type, error) {
				return Signed8(f), nil
			},
		),
		0xd1: bytes.Cast(
			bytes.ReadAs[int16](2, "int16", binary.BigEndian),
			func(f int16) (Type, error) {
				return Signed16(f), nil
			},
		),
		0xd2: bytes.Cast(
			bytes.ReadAs[int32](4, "int32", binary.BigEndian),
			func(f int32) (Type, error) {
				return Signed32(f), nil
			},
		),
		0xd3: bytes.Cast(
			bytes.ReadAs[int64](8, "int64", binary.BigEndian),
			func(f int64) (Type, error) {
				return Signed64(f), nil
			},
		),

		// fixext
		0xd4: extParser("fixext(1)", bytes.Const[uint8](1)),
		0xd5: extParser("fixext(2)", bytes.Const[uint8](2)),
		0xd6: extParser("fixext(4)", bytes.Const[uint8](4)),
		0xd7: extParser("fixext(8)", bytes.Const[uint8](8)),
		0xd8: extParser("fixext(16)", bytes.Const[uint8](16)),

		// strings
		0xd9: stringParser(bytes.ReadAs[uint8](1, "string", binary.BigEndian)),
		0xda: stringParser(bytes.ReadAs[uint16](2, "string", binary.BigEndian)),
		0xdb: stringParser(bytes.ReadAs[uint32](4, "string", binary.BigEndian)),
	}

	// positive fixint
	for i := byte(0x00); i <= byte(0x7f); i++ {
		cases[i] = bytes.Const[Type](Unsigned8(i))
	}

	// fixstring parser
	for i := byte(0xa0); i <= byte(0xbf); i++ {
		cases[i] = stringParser(
			bytes.Const[int](int(i - 0xa0)),
		)
	}

	// negative fixint
	for i, z := byte(0xe0), -32; i <= byte(0xff); i, z = i+1, z+1 {
		cases[i] = bytes.Const[Type](Signed8(z))
		if i == 0xff { // avoid endless loop
			break
		}
	}

	valuesParser := common.MapAs[byte, int, byte, Type](
		"expected positive fixint, fixstring or negative fixint",
		cases,
		bytes.Any(),
	)

	// fixarray parser
	for i := byte(0x90); i <= byte(0x9f); i++ {
		cases[i] = arrayParser(
			bytes.Const[int](int(i-0x90)),
			valuesParser,
		)
	}

	// array 16
	cases[0xdc] = arrayParser(
		bytes.ReadAs[uint16](2, "", binary.BigEndian),
		valuesParser,
	)

	// array 32
	cases[0xdd] = arrayParser(
		bytes.ReadAs[uint32](4, "expected array32", binary.BigEndian),
		valuesParser,
	)

	// fixmap parser
	for i := byte(0x80); i <= byte(0x8f); i++ {
		cases[i] = mapParser(
			bytes.Const[int](int(i-0x80)),
			valuesParser,
		)
	}

	// map 16 parser
	cases[0xde] = mapParser(
		bytes.ReadAs[uint16](2, "expected map16", binary.BigEndian),
		valuesParser,
	)

	// map 32 parser
	cases[0xdf] = mapParser(
		bytes.ReadAs[uint32](4, "expected map32", binary.BigEndian),
		valuesParser,
	)

	return valuesParser
}

func stringParser[T bytes.Number](
	parseSize common.Combinator[byte, int, T],
) common.Combinator[byte, int, Type] {
	return func(buffer common.Buffer[byte, int]) (Type, common.Error[int]) {
		size, err := parseSize(buffer)
		if err != nil {
			return String(""), err
		}

		pos := buffer.Position()

		data := make([]byte, int(size))
		for i := 0; i < int(size); i++ {
			x, err := buffer.Read(true)
			if err != nil {
				return String(""), common.NewParseError(
					pos,
					fmt.Sprintf("expected %d bytes of string data", int(size)),
				)
			}

			data[i] = x
		}

		return String(string(data)), nil
	}
}

func binaryParser[T bytes.Number](size int) common.Combinator[byte, int, Type] {
	comb := bytes.ReadAs[T](
		size,
		fmt.Sprintf("expected %d bytes of 'binary' data type", size),
		binary.BigEndian,
	)

	return func(buffer common.Buffer[byte, int]) (Type, common.Error[int]) {
		size, err := comb(buffer)
		if err != nil {
			return nil, err
		}

		pos := buffer.Position()

		data := make([]byte, int(size))
		for i := 0; i < int(size); i++ {
			x, err := buffer.Read(true)
			if err != nil {
				return nil, common.NewParseError(
					pos,
					fmt.Sprintf("expecte %d bytes of `binary` data type", int(size)),
				)
			}

			data[i] = x
		}

		return Binary(data), nil
	}
}

func arrayParser[T bytes.Number](
	parseSize common.Combinator[byte, int, T],
	parseValue common.Combinator[byte, int, Type],
) common.Combinator[byte, int, Type] {
	return func(buffer common.Buffer[byte, int]) (Type, common.Error[int]) {
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

func mapParser[T bytes.Number](
	parseSize common.Combinator[byte, int, T],
	parseValue common.Combinator[byte, int, Type],
) common.Combinator[byte, int, Type] {
	return func(buffer common.Buffer[byte, int]) (Type, common.Error[int]) {
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

func extParser[T bytes.Number](
	errMessage string,
	parseSize common.Combinator[byte, int, T],
) common.Combinator[byte, int, Type] {
	nameParser := bytes.ReadAs[int8](1, "name of ext type", binary.BigEndian)

	return func(buffer common.Buffer[byte, int]) (Type, common.Error[int]) {
		pos := buffer.Position()

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
				return nil, common.NewParseError(pos, errMessage)
			}

			data[i] = x
		}

		return Ext{Name: name, Data: data}, nil
	}
}
