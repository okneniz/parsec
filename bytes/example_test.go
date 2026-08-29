package bytes_test

import (
	"encoding/binary"
	"fmt"

	. "github.com/okneniz/parsec/bytes"
)

// ReadAs decodes a fixed number of bytes as a single number
// in the given byte order.
func ExampleReadAs() {
	be := ReadAs[uint16](2, "expected uint16", binary.BigEndian)
	le := ReadAs[uint16](2, "expected uint16", binary.LittleEndian)

	data := []byte{0x01, 0x02}

	result, err := Parse(data, be)
	fmt.Println(result, err)

	result, err = Parse(data, le)
	fmt.Println(result, err)

	// Output:
	// 258 <nil>
	// 513 <nil>
}

// SequenceOf expects a fixed sequence of bytes, like a magic number.
func ExampleSequenceOf() {
	magic := SequenceOf("expected magic bytes", 0xCA, 0xFE)

	result, err := Parse([]byte{0xCA, 0xFE, 0x00}, magic)
	fmt.Printf("% x %v\n", result, err)
	// Output: ca fe <nil>
}

// Range matches a single byte inside an inclusive range.
func ExampleRange() {
	printableASCII := Range("expected printable ASCII", 0x20, 0x7E)

	text := Many(0, Try(printableASCII))

	result, err := Parse([]byte("hi!"), text)
	fmt.Printf("%q %v\n", string(result), err)
	// Output: "hi!" <nil>
}
