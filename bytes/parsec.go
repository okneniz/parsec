package bytes

import (
	"github.com/okneniz/parsec"
)

// Parse applies the c combinator to data.
func Parse[T any](data []byte, parse parsec.Combinator[byte, int, T]) (T, error) {
	buf := Buffer(data)
	return parsec.Parse[byte, int, T](buf, parse)
}

// ParseString converts str to bytes and applies the c combinator to them.
func ParseString[T any](str string, parse parsec.Combinator[byte, int, T]) (T, error) {
	return Parse([]byte(str), parse)
}

// ParseFile reads the file and applies the c combinator to its content.
func ParseFile[T any](path string, parse parsec.Combinator[byte, int, T]) (T, error) {
	buf, err := BufferFromFile(path)
	if err != nil {
		var t T
		return t, err
	}

	return parsec.Parse[byte, int, T](buf, parse)
}
