package png

import (
	"testing"

	b "github.com/okneniz/parsec/bytes"
	. "github.com/okneniz/parsec/testing"
)

func TestPNG(t *testing.T) {
	comb := PNG()
	result, err := b.ParseFile("nibbler.png", comb)
	Check(t, err)
	t.Log("\n", result)
}
