package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Parse applies the c combinator to data.
func Parse[T any](data []byte, parse common.Combinator[byte, int, T]) (T, error) {
	buf := Buffer(data)
	return common.Parse[byte, int, T](buf, parse)
}

// ParseFile reads the file and applies the c combinator to its content.
func ParseFile[T any](path string, parse common.Combinator[byte, int, T]) (T, error) {
	buf, err := BufferFromFile(path)
	if err != nil {
		var t T
		return t, err
	}

	return common.Parse[byte, int, T](buf, parse)
}
