package strings

import (
	"fmt"
)

// Position - position in text.
type Position struct {
	line   uint
	column uint
	index  int
}

// Line - line number.
func (p Position) Line() uint {
	return p.line
}

// Column - column number.
func (p Position) Column() uint {
	return p.column
}

// String - return string representation of opsition.
func (p Position) String() string {
	return fmt.Sprintf("line=%d column=%d index=%d", p.line, p.column, p.index)
}
