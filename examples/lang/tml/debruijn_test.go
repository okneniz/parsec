package tml

import (
	"strings"
	"testing"
)

// exprOf parses src as val _ = <expr> and returns the expression.
func exprOf(t *testing.T, src string) Expr {
	t.Helper()

	decls, err := Parse("val _ = " + src)
	if err != nil {
		t.Fatalf("Parse of %q failed: %v", src, err)
	}

	return decls[0].(ValDecl).Binds[0].E
}

func TestDeBruijnIndices(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		// the innermost binder of the same name wins
		{"fn x => x", "(fn var(0))"},
		{"fn x => fn x => x", "(fn (fn var(0)))"},
		{"fn x => fn y => x", "(fn (fn var(1)))"},
		{"fn x => fn y => y", "(fn (fn var(0)))"},

		// a tuple pattern binds right to left: b innermost
		{"fn (a, b) => (a, b)", "(fn (tuple var(1) var(0)))"},

		// a let binding scopes over its body and the sibling decls
		{"let val a = 1 val b = 2 in (a, b) end", "(let (val (a 1)) (val (b 2)) in (tuple var(1) var(0)) end)"},
		{"fn x => let val y = 1 in x + y end", "(fn (let (val (y 1)) in (+ var(1) var(0)) end))"},

		// the builtin ~ sits under the initial stack
		{"fn x => ~x", "(fn (app var(1) var(0)))"},

		// a recursive name binds inside its own body
		{"let fun f x = f x in f end", "(let (fun f x (app var(1) var(0))) in var(0) end)"},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			got := exprOf(t, test.src)

			if got.String() != test.want {
				t.Errorf("expression %q = %s, want %s", test.src, got, test.want)
			}
		})
	}
}

func TestDeBruijnMutualRecursion(t *testing.T) {
	// the forward reference names a fun declared later: the fixup
	// pass turns it into a slot of the shared top-level frame, while
	// the already-declared name reads as an ordinary index
	decls, err := Parse(`
fun even n = if n = 0 then true else odd (n - 1)
fun odd n = if n = 0 then false else even (n - 1)
`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	want := []string{
		"(fun even n (if (= var(0) 0) true (app top(2) (- var(0) 1))))",
		"(fun odd n (if (= var(0) 0) false (app var(2) (- var(0) 1))))",
	}

	for i, d := range decls {
		if d.String() != want[i] {
			t.Errorf("declaration %d:\n got %s\nwant %s", i, d, want[i])
		}
	}
}

func TestDeBruijnUnbound(t *testing.T) {
	// a name that nothing binds survives the fixup pass and the
	// typechecker reports it; a later val does not count — only the
	// fun names are pre-bound
	tests := []struct {
		src  string
		want string
	}{
		{"val x = nosuch", "unbound variable nosuch"},
		{"val a = b\nval b = 1", "unbound variable b"},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
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
