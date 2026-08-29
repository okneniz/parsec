// This file holds the runnable documentation examples: the parse
// tree of a program, the inferred types of its bindings, and the
// values of its evaluation. They appear on pkg.go.dev as the package
// examples.
package tml

import (
	"fmt"
)

func Example() {
	decls, err := Parse(`
fun even n = if n = 0 then true else odd (n - 1)
fun odd n = if n = 0 then false else even (n - 1)

val _ = let val x = even 10 in if x then 1 else ~1 end
`)
	if err != nil {
		panic(err)
	}

	for _, d := range decls {
		fmt.Println(d)
	}
	// Output:
	// (fun even n (if (= n 0) true (app odd (- n 1))))
	// (fun odd n (if (= n 0) false (app even (- n 1))))
	// (val (_ (let (val (x (app even 10))) in (if x 1 (app ~ 1)) end)))
}

func Example_infer() {
	types, err := Infer(`
val id = fn x => x
val compose = fn f => fn g => fn x => f (g x)

fun add x y = x + y
`)
	if err != nil {
		panic(err)
	}

	for _, name := range []string{"id", "compose", "add"} {
		fmt.Printf("%s : %s\n", name, types[name])
	}
	// Output:
	// id : 'a -> 'a
	// compose : ('a -> 'b) -> ('c -> 'a) -> 'c -> 'b
	// add : int -> int -> int
}

func Example_eval() {
	env, err := Eval(`
fun fact n = if n = 0 then 1 else n * fact (n - 1)
fun fib n = if n < 2 then n else fib (n - 1) + fib (n - 2)

val answer = fact 5 + fib 10
val pair = (answer, answer < 100)
`)
	if err != nil {
		panic(err)
	}

	fmt.Println(env)
	// Output:
	// answer = 175
	// fact = <fn>
	// fib = <fn>
	// pair = (175, false)
}
