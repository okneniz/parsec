package common

type Buffer[T any, S any] interface {
	Read(bool) (T, error)
	Seek(S)
	Position() S
	IsEOF() bool
}
