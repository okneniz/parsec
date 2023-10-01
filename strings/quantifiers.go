package strings

import (
	p "github.com/okneniz/parsec/common"
)

// Optional - use c combinator to consume input data from buffer.
// If it failed, than return def value.
func Optional[T any](
	c p.Combinator[rune, Position, T],
	def T,
) p.Combinator[rune, Position, T] {
	return p.Optional[rune, Position, T](c, def)
}

// Many - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an empty slice even if nothing could be parsed.
func Many[T any](
	cap int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Many[rune, Position, T](cap, c)
}

// Some - accumulate data which returned by c consumer until it possible.
// Stop on first error or end of buffer.
// Returns an error if at least one element could not be read.
func Some[T any](
	cap int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Some[rune, Position, T](cap, c)
}

// Count - try to read X item by c combinator.
// Stop on first error.
func Count[T any](
	n int,
	c p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, []T] {
	return p.Count[rune, Position, T](n, c)
}
