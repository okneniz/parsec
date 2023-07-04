package strings

import (
	"fmt"
)

type Position struct {
	line   uint
	column uint
	index  int
}

func (p Position) Line() uint {
	return p.line
}

func (p Position) Column() uint {
	return p.column
}

func (p Position) String() string {
	return fmt.Sprintf("line=%d column=%d index=%d", p.line, p.column, p.index)
}
