package message_pack

import (
	"fmt"
	"reflect"
)

type Type interface {
	Type() string
	String() string
}

// Integer represents an integer

type Signed8 int8

func (x Signed8) Type() string {
	return "int8"
}

func (x Signed8) String() string {
	return fmt.Sprintf("%d", x)
}

type Signed16 int16

func (x Signed16) Type() string {
	return "int16"
}

func (x Signed16) String() string {
	return fmt.Sprintf("%d", x)
}

type Signed32 int32

func (x Signed32) Type() string {
	return "int8"
}

func (x Signed32) String() string {
	return fmt.Sprintf("%d", x)
}

type Signed64 int64

func (x Signed64) Type() string {
	return "int8"
}

func (x Signed64) String() string {
	return fmt.Sprintf("%d", x)
}

type Unsigned8 uint8

func (x Unsigned8) Type() string {
	return "uint8"
}

func (x Unsigned8) String() string {
	return fmt.Sprintf("%d", x)
}

type Unsigned16 uint16

func (x Unsigned16) Type() string {
	return "uint16"
}

func (x Unsigned16) String() string {
	return fmt.Sprintf("%d", x)
}

type Unsigned32 uint32

func (x Unsigned32) Type() string {
	return "uint32"
}

func (x Unsigned32) String() string {
	return fmt.Sprintf("%d", x)
}

type Unsigned64 uint64

func (x Unsigned64) Type() string {
	return "uint8"
}

func (x Unsigned64) String() string {
	return fmt.Sprintf("%d", x)
}

// Nil represents nil

type Nil struct{}

func (x Nil) Type() string {
	return "nil"
}

func (x Nil) String() string {
	return "nil"
}

// Boolean represents true or false

type Boolean bool

func (x Boolean) Type() string {
	return "bool"
}

func (x Boolean) String() string {
	return fmt.Sprintf("%t", x)
}

// Float represents a IEEE 754 double precision floating point number including NaN and Infinity

type Float32 float32

func (x Float32) Type() string {
	return "float32"
}

func (x Float32) String() string {
	return fmt.Sprintf("%f", x)
}

type Float64 float64

func (x Float64) Type() string {
	return "float64"
}

func (x Float64) String() string {
	return fmt.Sprintf("%f", x)
}

// Raw

// String extending Raw type represents a UTF-8 string

type String string

func (x String) Type() string {
	return "string"
}

func (x String) String() string {
	return fmt.Sprintf("%s", string(x))
}

// Binary extending Raw type represents a byte array

type Binary []byte

func (x Binary) Type() string {
	return "binary"
}

func (x Binary) String() string {
	return fmt.Sprintf("%#v", x)
}

// Array represents a sequence of objects

type Array[T any] []T

func (x Array[T]) Type() string {
	return "array"
}

func (x Array[T]) String() string {
	return fmt.Sprintf("%v", []T(x))
}

// Map represents key-value pairs of objects

type Map[K comparable, V any] map[K]V

func (x Map[K,V]) Type() string {
	t := reflect.TypeOf(x)
	return fmt.Sprintf("Map[%T,%T]", t.Key().Name(), t.Elem())
}

func (x Map[K,V]) String() string {
	return fmt.Sprintf("%#v", map[K]V(x))
}
