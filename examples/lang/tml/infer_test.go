package tml

import (
	"strings"
	"testing"
)

func TestInfer(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want map[string]string
	}{
		{
			name: "literals",
			src:  "val one = 1\nval yes = true",
			want: map[string]string{"one": "int", "yes": "bool"},
		},
		{
			name: "identity is polymorphic",
			src:  "val id = fn x => x",
			want: map[string]string{"id": "'a -> 'a"},
		},
		{
			name: "tuples and projections",
			src:  "val f = fn (x, y) => x",
			want: map[string]string{"f": "'a * 'b -> 'a"},
		},
		{
			name: "conditionals",
			src:  "val f = fn p => if p then 1 else 0",
			want: map[string]string{"f": "bool -> int"},
		},
		{
			name: "let-polymorphism instantiates",
			src: `
				val id = fn x => x
				val pair = (id 1, id true)
			`,
			want: map[string]string{"id": "'a -> 'a", "pair": "int * bool"},
		},
		{
			name: "curried functions",
			src:  "fun add x y = x + y",
			want: map[string]string{"add": "int -> int -> int"},
		},
		{
			name: "tuple argument functions",
			src:  "fun add (a, b) = a + b",
			want: map[string]string{"add": "int * int -> int"},
		},
		{
			name: "mutual recursion at the top level",
			src: `
				fun even n = if n = 0 then true else odd (n - 1)
				fun odd n = if n = 0 then false else even (n - 1)
			`,
			want: map[string]string{"even": "int -> bool", "odd": "int -> bool"},
		},
		{
			name: "function composition",
			src:  "val compose = fn f => fn g => fn x => f (g x)",
			want: map[string]string{"compose": "('a -> 'b) -> ('c -> 'a) -> 'c -> 'b"},
		},
		{
			name: "inner let does not leak",
			src: `
				val x = let val y = fn z => z in y end
				val w = (x 1, x true)
			`,
			want: map[string]string{"x": "'a -> 'a", "w": "int * bool"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := inferOf(t, test.src)

			if len(got) != len(test.want) {
				t.Fatalf("Infer gave %d bindings (%v), want %d", len(got), got, len(test.want))
			}

			for name, want := range test.want {
				if got[name] != want {
					t.Errorf("type of %s = %s, want %s", name, got[name], want)
				}
			}
		})
	}
}

func TestInferErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"int and bool", "val x = 1 + true", "cannot unify bool with int"},
		{"condition not boolean", "val x = if 1 then 2 else 3", "cannot unify"},
		{"branches disagree", "val x = if true then 1 else false", "cannot unify"},
		{"unbound variable", "val x = y", "unbound variable y"},
		{"infinite type", "val x = fn x => x x", "infinite type"},
		{"operator on tuple", "val x = (1, 2) + 3", "cannot unify"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Infer(test.src)
			if err == nil {
				t.Fatalf("Infer(%q) succeeded, want error", test.src)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("Infer(%q) error = %v, want it to contain %q", test.src, err, test.want)
			}
		})
	}
}

// inferOf typechecks src and returns the printed type of the binding.
func inferOf(t *testing.T, src string) map[string]string {
	t.Helper()

	types, err := Infer(src)
	if err != nil {
		t.Fatalf("Infer failed: %v", err)
	}

	out := map[string]string{}

	for name, ty := range types {
		out[name] = ty.String()
	}

	return out
}
