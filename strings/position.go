package strings

import (
	"fmt"
)

// Position is a location in text: line and column are zero-based,
// index is the absolute rune offset from the beginning.
type Position struct {
	line   uint
	column uint
	index  int
}

// Line returns the zero-based line number.
func (p Position) Line() uint {
	return p.line
}

// Column returns the zero-based column number.
func (p Position) Column() uint {
	return p.column
}

// String returns a human-readable representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("line=%d column=%d index=%d", p.line, p.column, p.index)
}
