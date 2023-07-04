package testing

import (
	"testing"
)

func TestDigitsToNum(t *testing.T) {
	examples := map[string]int{
		"0":     0,
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"10":    10,
		"11":    11,
		"101":   101,
		"123":   123,
		"10723": 10723,
		"50221": 50221,
	}

	for input, expected := range examples {
		actual, err := DigitsToNum([]rune(input))
		if err != nil {
			t.Error(err)
			t.Errorf("expected %v, actual %v - input %v", expected, actual, input)
			continue
		}
		if expected != actual {
			t.Errorf("expected %v, actual %v - input %v", expected, actual, input)
		}
	}
}
