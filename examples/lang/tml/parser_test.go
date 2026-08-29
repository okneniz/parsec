package tml

import (
	"strings"
	"testing"
)

func TestPrecedence(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"42", "42"},
		{"true", "true"},
		{"1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"a - b - c", "(- (- a b) c)"},
		{"a = b andalso c < d", "(andalso (= a b) (< c d))"},
		{"a orelse b andalso c", "(orelse a (andalso b c))"},
		{"f x y", "(app (app f x) y)"},
		{"f ~3", "(app f (app ~ 3))"},
		{"(1 + 2) * 3", "(* (+ 1 2) 3)"},
		{"if p then f x else g y", "(if p (app f x) (app g y))"},
		{"fn x => x + 1", "(fn x (+ x 1))"},
		{"fn (x, _) => x", "(fn (tuple x _) x)"},
		{"(a, b, c)", "(tuple a b c)"},
		{"let val x = 1 in x end", "(let (val (x 1)) in x end)"},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			decls, err := Parse("val _ = " + test.src)
			if err != nil {
				t.Fatalf("Parse of %q failed: %v", test.src, err)
			}

			bind := decls[0].(ValDecl).Binds[0]

			if bind.E.String() != test.want {
				t.Errorf("expression %q = %s, want %s", test.src, bind.E, test.want)
			}
		})
	}
}

func TestProgram(t *testing.T) {
	src := `
(* factorial and fibonacci *)
fun fact n = if n = 0 then 1 else n * fact (n - 1)

fun fib n = if n < 2 then n else fib (n - 1) + fib (n - 2)

fun add (a, b) = a + b

val answer = let
  val five = fact 5
  val ten = fib 10
  val seven = add (3, 4)
in
  five + ten + seven - ~1
end
`

	want := []string{
		"(fun fact n (if (= n 0) 1 (* n (app fact (- n 1)))))",
		"(fun fib n (if (< n 2) n (+ (app fib (- n 1)) (app fib (- n 2)))))",
		"(fun add (tuple a b) (+ a b))",
		"(val (answer (let (val (five (app fact 5))) (val (ten (app fib 10))) (val (seven (app add (tuple 3 4)))) in (- (+ (+ five ten) seven) (app ~ 1)) end)))",
	}

	decls, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(decls) != len(want) {
		t.Fatalf("parsed %d declarations, want %d", len(decls), len(want))
	}

	for i, d := range decls {
		if d.String() != want[i] {
			t.Errorf("declaration %d:\n got %s\nwant %s", i, d, want[i])
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"missing then", "val x = if true 1 else 2"},
		{"reserved word as expression", "val x = val"},
		{"dangling operator", "val x = 1 +"},
		{"unclosed paren", "val x = (1 + 2"},
		{"unknown declaration", "datatype t = A"},
		{"fn without arrow", "val f = fn x x"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.src)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want error", test.src)
			}

			// the declaration list stops at the first declaration
			// that does not parse, so every broken program is
			// reported by the end check
			if !strings.Contains(err.Error(), "expected declaration") {
				t.Errorf("Parse(%q) error = %v, want it to contain %q", test.src, err, "expected declaration")
			}
		})
	}
}
