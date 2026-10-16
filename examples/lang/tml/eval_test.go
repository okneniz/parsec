package tml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// evalOf interprets src and returns the value of its main.
func evalOf(t *testing.T, src string) Value {
	t.Helper()

	value, err := Eval(src)
	require.NoError(t, err)

	return value
}

func TestEval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"factorial", "fun fact n = if n = 0 then 1 else n * fact (n - 1)\nfun main () = fact 5", "120"},
		{"fibonacci", "fun fib n = if n < 2 then n else fib (n - 1) + fib (n - 2)\nfun main () = fib 10", "55"},
		{"negation", "fun main () = ~3 + 10", "7"},
		{"tuples", "fun main () = (1, true, (2, false))", "(1, true, (2, false))"},
		{"unit", "fun main () = ()", "()"},
		{"tuple pattern", "fun main () = let val (a, b) = (7, false) in b end", "false"},
		{"tuple argument", "fun add (a, b) = a + b\nfun main () = add (3, 4)", "7"},
		{"shadowing", "fun main () = let val x = 1 in let val x = true in x end end", "true"},
		{"higher order", "fun twice f x = f (f x)\nfun main () = twice (fn n => n + 3) 4", "10"},
		{"polymorphic application", "fun main () = let val id = fn x => x in (id 1, id true) end", "(1, true)"},
		{
			"mutual recursion",
			"fun even n = if n = 0 then true else odd (n - 1)\nfun odd n = if n = 0 then false else even (n - 1)\nfun main () = even 10",
			"true",
		},
		{"local fun recursion", "fun main () = let fun f x = if x = 0 then 0 else f (x - 1) in f 5 end", "0"},
		{"escaped closure", "fun main () = let val f = fn y => 10 + y in f 7 end", "17"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := renderValue(evalOf(t, test.src))

			require.Equal(t, test.want, got)
		})
	}
}

func TestEvalErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"division by zero", "fun main () = 1 / 0", "division by zero"},
		{"match failure", "fun main () = let val 1 = 2 in () end", "pattern match failure"},
		{"type error first", "fun main () = 1 + true", "cannot unify"},
		{"no main", "val x = 1", "no main function"},
		{"unbound name", "fun main () = nosuch", "unbound variable"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := Eval(test.src)
			require.ErrorContains(t, err, test.want)
		})
	}
}
