package common

// Buffer - input data for combinators.
// T - type of input data
// P - type of position
type Buffer[T any, P any] interface {
	// Read - read next item.
	// If greedy is true, the buffer advances and keeps the new position,
	// including the case when the calling combinator later fails on this item;
	// with greedy false the item is only peeked and the position stays unchanged.
	Read(greedy bool) (T, error)
	// Seek - change buffer position
	Seek(position P) error
	// Position - return current buffer position
	Position() P
	// IsEOF - true if buffer ended.
	IsEOF() bool
}
