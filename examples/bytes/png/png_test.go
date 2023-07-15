package png

import (
	"testing"

	b "git.sr.ht/~okneniz/parsec/bytes"
	. "git.sr.ht/~okneniz/parsec/testing"
)

func TestPNG(t *testing.T) {
	comb := PNG()
	result, err := b.ParseFile("nibbler.png", comb)
	Check(t, err)
	t.Log("\n", result)
}
