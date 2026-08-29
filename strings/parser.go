package strings

import (
	"os"

	"github.com/okneniz/parsec"
)

// Parse applies the c combinator to data.
func Parse[T any](
	data []rune,
	parse parsec.Combinator[rune, Position, T],
) (T, parsec.Error[Position]) {
	buf := Buffer(data)
	return parsec.Parse[rune, Position, T](buf, parse)
}

// ParseString converts str to runes and applies the c combinator to them.
func ParseString[T any](
	str string,
	parse parsec.Combinator[rune, Position, T],
) (T, parsec.Error[Position]) {
	return Parse([]rune(str), parse)
}

// ParseFile reads the file, converts its content to runes and applies
// the c combinator to them.
func ParseFile[T any](
	path string,
	parse parsec.Combinator[rune, Position, T],
) (T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		var t T

		return t, err
	}

	return ParseString(string(data), parse)
}
