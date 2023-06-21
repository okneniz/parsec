package parsec

import (
	"fmt"
	"testing"
)

func TestChainl(t *testing.T) {
	next := Satisfy(true, Anything[byte])

	comb := Chainl(
		func(buffer Buffer[byte]) (string, bool) {
			x, ok := next(buffer)
			if !ok {
				return "", false
			}

			return string(x), true
		},
		func(buffer Buffer[byte]) (func(string, string) string, bool) {
			return func(x, y string) string {
				return fmt.Sprintf("(%v %v)", x, y)
			}, true
		},
		"default",
	)

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	next := Satisfy(true, Anything[byte])

	comb := Chainl1(
		func(buffer Buffer[byte]) (string, bool) {
			x, ok := next(buffer)
			if !ok {
				return "", false
			}

			return string(x), true
		},
		func(buffer Buffer[byte]) (func(string, string) string, bool) {
			return func(x, y string) string {
				return fmt.Sprintf("(%v %v)", x, y)
			}, true
		},
	)

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	next := Satisfy(true, Anything[byte])

	comb := Chainr(
		func(buffer Buffer[byte]) (string, bool) {
			x, ok := next(buffer)
			if !ok {
				return "", false
			}

			return string(x), true
		},
		func(buffer Buffer[byte]) (func(string, string) string, bool) {
			return func(x, y string) string {
				return fmt.Sprintf("(%v %v)", x, y)
			}, true
		},
		"default",
	)

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	next := Satisfy(true, Anything[byte])

	comb := Chainr1(
		func(buffer Buffer[byte]) (string, bool) {
			x, ok := next(buffer)
			if !ok {
				return "", false
			}

			return string(x), true
		},
		func(buffer Buffer[byte]) (func(string, string) string, bool) {
			return func(x, y string) string {
				return fmt.Sprintf("(%v %v)", x, y)
			}, true
		},
	)

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}
