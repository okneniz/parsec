package message_pack

type Type interface {
	Type() string
	String() string
}

// Integer represents an integer

type Signed8 int8
type Signed16 int16
type Signed32 int32
type Signed64 int64

type Unsigned8 uint8
type Unsigned16 uint16
type Unsigned32 uint32
type Unsigned64 uint64

// Nil represents nil

type Nil struct{}

// Boolean represents true or false

type Boolean bool

// Float represents a IEEE 754 double precision floating point number including NaN and Infinity

type Float32 float32
type Float64 float64

// Raw

// String extending Raw type represents a UTF-8 string

type String string

// Binary extending Raw type represents a byte array

type Binary []byte

// Array represents a sequence of objects

type Array[T any] []T

// Map represents key-value pairs of objects

type Map[K comparable, V any] map[K]V
