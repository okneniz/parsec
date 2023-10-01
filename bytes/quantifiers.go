package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Optional - use c combinator to consume input byte from buffer.
// If it failed, than return def value.
func Optional[T any](
	c p.Combinator[byte, int, T],
	def T,
) p.Combinator[byte, int, T] {
	return p.Optional[byte, int, T](c, def)
}

// Many - read bytes and accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an empty slice even if nothing could be parsed.
func Many[T any](
	cap int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Many[byte, int, T](cap, c)
}

// Some - read bytes and accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an error if at least one element could not be read.
func Some[T any](
	cap int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Some[byte, int, T](cap, c)
}

// Count - try to read X item by c combinator.
// Stop on first error.
func Count[T any](
	n int,
	c p.Combinator[byte, int, T],
) p.Combinator[byte, int, []T] {
	return p.Count[byte, int, T](n, c)
}
