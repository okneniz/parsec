// Package tml is a parser for tiny ml: a minimal Standard ML dialect
// small enough to read in one sitting. The language has exactly what
// a first functional language needs — let bindings, recursive
// functions, if-then-else, integers, booleans and the unit value () —
// and nothing else: no modules, no datatypes, no references, no loops.
//
// A program is a sequence of val and fun declarations, plus the main
// function Eval calls with the unit value; fun bindings are recursive
// by default:
//
//	fun fact n = if n = 0 then 1 else n * fact (n - 1)
//	fun main () = let val five = fact 5 in five + ~1 end
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
// index — 0 is the innermost binder. The top level binds in two
// passes: the first collects the names of the functions, the second —
// with the names already on the stack — resolves every reference to
// them, so the functions see each other regardless of declaration
// order: that is the mutual recursion of tiny ml. A name no binder
// holds is a FreeVar, and the typechecker reports it.
//
// On top of the parser sit the two classic semantic passes, both
// addressing the binder row positionally:
//
//   - Infer, a Hindley-Milner typechecker with let-polymorphism: a
//     mutation-based unification with occurs check, generalization
//     against the environment, and instantiation at use sites. Top
//     level functions are pre-bound with shared type variables, so
//     they may be mutually recursive;
//   - Eval, a substitution-based interpreter running after the
//     typecheck. The beta step is literal: apply matches the argument
//     pattern and moves the matched values into the body through
//     substitutionTops — the de Bruijn shift/subst machinery of
//     subst.go — and the vals of the top level and of every let leave
//     the terms the same way, so a running term carries no val slots.
//     What substitution cannot tie is the recursive knot: the
//     functions of a scope — the top level included — become one row
//     of knots, each holding itself in its own slot, and a free index
//     of a running term always names a knot. Eval calls main with ()
//     and returns its value. Division by zero and pattern match
//     failures are the only runtime errors left. EvalFile reads the
//     program from a file.
//
// The programs directory holds ten teaching programs — from literals
// to mutual recursion — with the language features explained in
// comments; programs_test.go runs them all and checks the value each
// main computes.
package tml
