package bytes

import (
	"github.com/okneniz/parsec"
)

// Parse applies the c combinator to data with the given state.
func Parse[O any, S any](state S, data []byte, parse parsec.Combinator[byte, int, O, S]) (O, error) {
	buf := Buffer(data)
	return parsec.Parse[byte, int, O, S](state, buf, parse)
}

// ParseString converts str to bytes and applies the c combinator
// to them with the given state.
func ParseString[O any, S any](state S, str string, parse parsec.Combinator[byte, int, O, S]) (O, error) {
	return Parse(state, []byte(str), parse)
}

// ParseFile reads the file and applies the c combinator to its content
// with the given state.
func ParseFile[O any, S any](state S, path string, parse parsec.Combinator[byte, int, O, S]) (O, error) {
	buf, err := BufferFromFile(path)
	if err != nil {
		var o O
		return o, err
	}

	return parsec.Parse[byte, int, O, S](state, buf, parse)
}
