package parsec

// Buffer is the input of every [Combinator]: a readable, seekable
// sequence of items of type T with positions of type P.
// The standard implementations are
// [github.com/okneniz/parsec/strings.Buffer] for runes and
// [github.com/okneniz/parsec/bytes.Buffer] for bytes.
type Buffer[T any, P any] interface {
	// Read returns the next item.
	// If greedy is true, the buffer advances and keeps the new position,
	// including the case when the calling combinator later fails on this item;
	// with greedy false the item is only peeked and the position stays unchanged.
	Read(greedy bool) (T, error)

	// Seek moves the buffer to a previously obtained position.
	// It returns [ErrOutOfBounds] when the position is invalid.
	Seek(position P) error

	// Position returns the current buffer position.
	Position() P

	// IsEOF reports whether the buffer is fully consumed.
	IsEOF() bool
}
