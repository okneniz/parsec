package parsec

import (
	"testing"
)

func TestCombnators(t *testing.T) {
	// a := byte('a')
	// b := byte('b')
	c := byte('c')

	comb := Satisfy[byte](true, func(x byte) bool { return x != c })

	result, ok := ParseBytes([]byte("a"), comb)
	t.Logf("result %v, ok %v", result, ok)

	result, ok = ParseBytes([]byte("b"), comb)
	t.Logf("result %v, ok %v", result, ok)

	result, ok = ParseBytes([]byte("c"), comb)
	t.Logf("result %v, ok %v", result, ok)
}
