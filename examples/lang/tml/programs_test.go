package tml

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPrograms runs every teaching program from the programs
// directory and checks the value its main computes, in declaration
// order. The programs double as documentation: each one explains a
// language feature in comments.
func TestPrograms(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"01-basics.tml", "(42, 7, 9, 5, true, true)"},
		{"02-functions.tml", "(<fn>, 6, <fn>, 10, <fn>, 120, <fn>, 55)"},
		{"03-higher-order.tml", "(<fn>, (1, true), <fn>, 10, <fn>, <fn>, <fn>, 7)"},
		{
			"04-patterns.tml",
			"((3, 4), 7, <fn>, 10, <fn>, 7, <fn>, 25, <fn>, (-1, 0, 1))",
		},
		{"05-mutual.tml", "(<fn>, <fn>, true, <fn>, 61)"},
		{
			"06-numbers.tml",
			"(<fn>, <fn>, 21, <fn>, <fn>, (1024, 1024), <fn>, (true, false), <fn>, 31, 32)",
		},
		{"07-hofstadter.tml", "(<fn>, <fn>, (6, 6), <fn>, 6)"},
		{
			"08-rationals.tml",
			"(<fn>, <fn>, <fn>, <fn>, <fn>, (1, 2), (1, 3), (5, 6), (1, 2), (-7, 10), (1, 4))",
		},
		{
			"09-church.tml",
			"(<fn>, <fn>, <fn>, <fn>, <fn>, 2, 1, 2, <fn>, <fn>, <fn>, <fn>, <fn>, <fn>, <fn>, <fn>, <fn>, 4, 7, 8)",
		},
		{
			"10-lists.tml",
			"(<fn>, <fn>, <fn>, <fn>, <fn>, 3, 2, <fn>, 6, <fn>, 6)",
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			value, err := EvalFile(filepath.Join("programs", test.file))
			require.NoError(t, err)

			require.Equal(t, test.want, renderValue(value))
		})
	}
}
