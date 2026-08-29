package tml

import (
	"strings"
	"testing"
)

// evalOf interprets src and returns the value of the binding.
func evalOf(t *testing.T, src, name string) Value {
	t.Helper()

	env, err := Eval(src)
	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	v, bound := env[name]
	if !bound {
		t.Fatalf("binding %s not found in %v", name, env)
	}

	return v
}

func TestEval(t *testing.T) {
	tests := []struct {
		name string
		src  string
		bind string
		want string
	}{
		{"factorial", "fun fact n = if n = 0 then 1 else n * fact (n - 1)\nval x = fact 5", "x", "120"},
		{"fibonacci", "fun fib n = if n < 2 then n else fib (n - 1) + fib (n - 2)\nval x = fib 10", "x", "55"},
		{"negation", "val x = ~3 + 10", "x", "7"},
		{"tuples", "val x = (1, true, (2, false))", "x", "(1, true, (2, false))"},
		{"tuple pattern", "val (a, b) = (7, false)", "b", "false"},
		{"tuple argument", "fun add (a, b) = a + b\nval x = add (3, 4)", "x", "7"},
		{"shadowing", "val x = let val x = 1 in let val x = true in x end end", "x", "true"},
		{
			"higher order",
			"fun twice f x = f (f x)\nval x = twice (fn n => n + 3) 4",
			"x", "10",
		},
		{
			"polymorphic application",
			"val id = fn x => x\nval x = (id 1, id true)",
			"x", "(1, true)",
		},
		{
			"mutual recursion",
			"fun even n = if n = 0 then true else odd (n - 1)\nfun odd n = if n = 0 then false else even (n - 1)\nval x = even 10",
			"x", "true",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := evalOf(t, test.src, test.bind)

			if got.String() != test.want {
				t.Errorf("Eval = %s, want %s", got, test.want)
			}
		})
	}
}

func TestEvalErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"division by zero", "val x = 1 / 0", "division by zero"},
		{"match failure", "val 1 = 2", "pattern match failure"},
		{"type error first", "val x = 1 + true", "cannot unify"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Eval(test.src)
			if err == nil {
				t.Fatalf("Eval(%q) succeeded, want error", test.src)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("Eval(%q) error = %v, want it to contain %q", test.src, err, test.want)
			}
		})
	}
}
