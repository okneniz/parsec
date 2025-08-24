package png

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okneniz/parsec/bytes"
)

func TestPNG(t *testing.T) {
	comb := PNG()
	result, err := bytes.ParseFile("nibbler.png", comb)
	assert.NoError(t, err)
	t.Log("\n", result)
}
