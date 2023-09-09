package common

// Buffer - input data for combinators.
// T - type of input data
// P - type of position
type Buffer[T any, P any] interface {
	// Read - read next item, if greedy buffer keep position after reading.
	Read(greedy bool) (T, error)
	// Seek - change buffer position
	Seek(position P)
	// Position - return current buffer position
	Position() P
	// IsEOF - true if buffer ended.
	IsEOF() bool
}
