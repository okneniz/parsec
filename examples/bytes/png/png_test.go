package png

import (
	"testing"

	"github.com/okneniz/parsec"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/parsec/bytes"
)

func TestPNG(t *testing.T) {
	comb := PNG()
	result, err := bytes.ParseFile(parsec.Stateless{}, "nibbler.png", comb)
	require.NoError(t, err)
	t.Log("\n", result)
}
