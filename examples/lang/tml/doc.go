// Package tml is a parser for tiny ml: a minimal Standard ML dialect
// small enough to read in one sitting. The language has exactly what
// a first functional language needs — let bindings, recursive
// functions, if-then-else, integers and booleans — and nothing else:
// no modules, no datatypes, no references, no loops.
//
// A program is a sequence of val and fun declarations; fun bindings
// are recursive by default:
//
//	fun fact n = if n = 0 then 1 else n * fact (n - 1)
//	val answer = let val five = fact 5 in five + ~1 end
//
// The package is a two-stage pipeline in miniature: a rune-level
// lexer built with github.com/okneniz/parsec/strings — nested
// comments and maximal-munch operators included — feeds a lazy
// github.com/okneniz/parsec/tokens buffer which lexes on demand, and
// the parser runs over it with the same generic core: a static fixity
// table and the shunting-yard resolution of infix chains. The rune
// position doubles as the token position, so the tokens themselves
// stay coordinate-free.
//
// On top of the parser sit the two classic semantic passes:
//
//   - Infer, a Hindley-Milner typechecker with let-polymorphism: a
//     mutation-based unification with occurs check, generalization
//     against the environment, and instantiation at use sites. Top
//     level functions are pre-bound with shared type variables, so
//     they may be mutually recursive;
//   - Eval, an environment-based interpreter running after the
//     typecheck: closures capture their environment by reference,
//     curried declarations expand into nested fn-expressions, and
//     patterns match at bindings and calls. Division by zero and
//     pattern match failures are the only runtime errors left.
//     EvalFile reads the program from a file.
//
// The programs directory holds five teaching programs — from literals
// to mutual recursion — with the language features explained in
// comments; programs_test.go runs them all and checks the results.
package tml
