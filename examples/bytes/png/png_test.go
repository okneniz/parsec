package png

import (
	"testing"

	"github.com/okneniz/parsec/bytes"
	. "github.com/okneniz/parsec/testing"
)

func TestPNG(t *testing.T) {
	comb := PNG()
	result, err := bytes.ParseFile("nibbler.png", comb)
	Check(t, err)
	t.Log("\n", result)
}
