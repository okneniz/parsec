package strings

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
