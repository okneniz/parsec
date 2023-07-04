package bytes

import (
	"git.sr.ht/~okneniz/parsec/common"
)

func Parse[T any](data []byte, parse common.Combinator[byte, int, T]) (T, error) {
	buf := Buffer(data)
	return common.Parse[byte, int, T](buf, parse)
}
