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
// github.com/okneniz/parsec/tokens buffer which lexes on demand. The
// expression ladder is github.com/okneniz/parsec/lang: tiny ml
// declares its operators and the reserved words andalso and orelse —
// infixes of their own precedence level — on a lang.Parser, and its
// forms with it: if and fn at the expression level, let, the
// negation ~ and the parenthesis at the operand level, each form
// parsing its subexpressions through the ladder it is part of.
// Declarations and patterns are the package's own recursive descent
// around the ladder. The rune position doubles as the token position,
// so the tokens themselves stay coordinate-free.
//
// The whole pipeline carries one parser state: the binder stack.
// Every fn pattern, every declaration and every let pushes its
// binders for the nested scopes and pops them on the way out, and an
// identifier read off the stack becomes a Var node with its de Bruijn
// index — 0 is the innermost binder. A name off the stack is a
// ForwardRef the fixup pass of Parse resolves against the whole
// program: one that names a top-level fun declared later — the mutual
// recursion of tiny ml — becomes a TopVar, a slot of the shared
// top-level frame; whatever survives names nothing at all.
//
// On top of the parser sit the two classic semantic passes, both
// running on the positional environment the parser built:
//
//   - Infer, a Hindley-Milner typechecker with let-polymorphism: a
//     mutation-based unification with occurs check, generalization
//     against the environment, and instantiation at use sites. Top
//     level functions are pre-bound with shared type variables, so
//     they may be mutually recursive;
//   - Eval, an environment-based interpreter running after the
//     typecheck: the top level is one frame of cells the
//     declarations fill as the program runs — what a closure built
//     earlier sees a function declared later through — while the
//     nested scopes extend snapshots of it, curried declarations
//     expand into nested fn-expressions, and patterns match at
//     bindings and calls. Division by zero and pattern match
//     failures are the only runtime errors left. EvalFile reads the
//     program from a file.
//
// The programs directory holds ten teaching programs — from literals
// to mutual recursion — with the language features explained in
// comments; programs_test.go runs them all and checks the results.
package tml
