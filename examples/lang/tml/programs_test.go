package tml

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestPrograms runs every teaching program from the programs
// directory and checks its final bindings. The programs double as
// documentation: each one explains a language feature in comments.
func TestPrograms(t *testing.T) {
	tests := []struct {
		file string
		want []string
	}{
		{
			file: "01-basics.tml",
			want: []string{
				"answer = 42",
				"comparison = true",
				"grouped = 9",
				"logic = true",
				"negative = 5",
				"precedence = 7",
			},
		},
		{
			file: "02-functions.tml",
			want: []string{
				"add = <fn>",
				"fact = <fn>",
				"fact5 = 120",
				"fib = <fn>",
				"fib10 = 55",
				"six = 6",
				"successor = <fn>",
				"ten = 10",
			},
		},
		{
			file: "03-higher-order.tml",
			want: []string{
				"compose = <fn>",
				"double = <fn>",
				"id = <fn>",
				"inc = <fn>",
				"pair = (1, true)",
				"seven = 7",
				"ten = 10",
				"twice = <fn>",
			},
		},
		{
			file: "04-patterns.tml",
			want: []string{
				"add = <fn>",
				"dist = <fn>",
				"first = 10",
				"fiveSquared = 25",
				"fst = <fn>",
				"sign = <fn>",
				"signs = (-1, 0, 1)",
				"sum = 7",
				"sum2 = 7",
				"x = 3",
				"y = 4",
			},
		},
		{
			file: "05-mutual.tml",
			want: []string{
				"ack = <fn>",
				"ack33 = 61",
				"even = <fn>",
				"isTenEven = true",
				"odd = <fn>",
			},
		},
		{
			file: "06-numbers.tml",
			want: []string{
				"fastpow = <fn>",
				"gcd = <fn>",
				"gcdDemo = 21",
				"isPrime = <fn>",
				"isqrt = <fn>",
				"mod = <fn>",
				"power = <fn>",
				"powers = (1024, 1024)",
				"primality = (true, false)",
				"root1024 = 32",
				"root997 = 31",
			},
		},
		{
			file: "07-hofstadter.tml",
			want: []string{
				"female = <fn>",
				"male = <fn>",
				"pair10 = (6, 6)",
				"q = <fn>",
				"q10 = 6",
			},
		},
		{
			file: "08-rationals.tml",
			want: []string{
				"addR = <fn>",
				"fiveSixths = (5, 6)",
				"gcd = <fn>",
				"half = (1, 2)",
				"mod = <fn>",
				"mulR = <fn>",
				"normalized = (-7, 10)",
				"quarter = (1, 4)",
				"rat = <fn>",
				"reduced = (1, 2)",
				"third = (1, 3)",
			},
		},
		{
			file: "09-church.tml",
			want: []string{
				"andDemo = 2",
				"c0 = <fn>",
				"c1 = <fn>",
				"c2 = <fn>",
				"c3 = <fn>",
				"c4 = <fn>",
				"cand = <fn>",
				"cfalse = <fn>",
				"cmult = <fn>",
				"cnot = <fn>",
				"cor = <fn>",
				"cplus = <fn>",
				"csucc = <fn>",
				"ctrue = <fn>",
				"eight = 8",
				"four = 4",
				"notTrue = 2",
				"orDemo = 1",
				"seven = 7",
				"toInt = <fn>",
			},
		},
		{
			file: "10-lists.tml",
			want: []string{
				"append = <fn>",
				"cons = <fn>",
				"lenOneTwo = 2",
				"length = <fn>",
				"map = <fn>",
				"nil = <fn>",
				"oneTwo = <fn>",
				"sum = <fn>",
				"sumAppended = 6",
				"sumDoubled = 6",
				"sumOneTwo = 3",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			env, err := EvalFile(filepath.Join("programs", test.file))
			if err != nil {
				t.Fatalf("EvalFile failed: %v", err)
			}

			want := strings.Join(test.want, "\n")

			if got := env.String(); got != want {
				t.Errorf("bindings of %s:\n got:\n%s\nwant:\n%s", test.file, got, want)
			}
		})
	}
}
