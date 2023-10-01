package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Parse - parse bytes from input slice by c combinator.
func Parse[T any](data []byte, parse common.Combinator[byte, int, T]) (T, error) {
	buf := Buffer(data)
	return common.Parse[byte, int, T](buf, parse)
}

// ParseFile - parse bytes from file by c combinator.
func ParseFile[T any](path string, parse common.Combinator[byte, int, T]) (T, error) {
	buf, err := BufferFromFile(path)
	if err != nil {
		var t T
		return t, err
	}

	return common.Parse[byte, int, T](buf, parse)
}
