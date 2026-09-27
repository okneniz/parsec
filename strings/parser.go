package strings

import (
	"os"

	"github.com/okneniz/parsec"
)

// Parse applies the c combinator to data with the given state.
func Parse[O any, S any](
	state S,
	data []rune,
	parse parsec.Combinator[rune, Position, O, S],
) (O, parsec.Error[Position]) {
	buf := Buffer(data)
	return parsec.Parse[rune, Position, O, S](state, buf, parse)
}

// ParseString converts str to runes and applies the c combinator
// to them with the given state.
func ParseString[O any, S any](
	state S,
	str string,
	parse parsec.Combinator[rune, Position, O, S],
) (O, parsec.Error[Position]) {
	return Parse(state, []rune(str), parse)
}

// ParseFile reads the file, converts its content to runes and applies
// the c combinator to them with the given state.
func ParseFile[O any, S any](
	state S,
	path string,
	parse parsec.Combinator[rune, Position, O, S],
) (O, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		var t O

		return t, err
	}

	return ParseString(state, string(data), parse)
}
