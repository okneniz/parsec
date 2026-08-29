package common

// Integer - constraint for any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float - constraint for any float type.
type Float interface {
	~float32 | ~float64
}
