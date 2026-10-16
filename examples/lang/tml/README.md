# tiny ml

tiny ml (tml) is a minimal Standard ML dialect — `let` bindings,
recursive functions, `if-then-else` over integers, booleans and the
unit value `()` — with a complete front end and semantics on ~2000
lines of Go:

```
source → lexer → tokens → parser → AST → typechecker → interpreter → values
```

The point of the package is not the language, which is deliberately
tiny, but the route: every stage is a readable, self-contained
implementation of a classic technique, built on the parsec combinators.

- `fun fact n = if n = 0 then 1 else n * fact (n - 1)` parses, gets
  the type `int -> int` without a single annotation, and runs to 120
  on `fact 5`;
- `val id = fn x => x` gets `'a -> 'a`, and `(id 1, id true)` gets
  `int * bool` — polymorphism is inferred, not declared;
- with no datatypes in the language, the `programs/` directory still
  builds lists — a list is encoded as its own fold.

A program's entry point is `fun main () = ...`: Eval calls it with
the unit value and returns what it computes.

## The map

Read the package in this order; every file is one stage of the
pipeline:

| file | stage | the classic technique inside |
| --- | --- | --- |
| `lexer.go` | text → tokens | maximal munch, nested comments; produces `tokens.Token` from the shared [tokens](../../tokens) layer |
| `parser.go` | tokens → AST | the expression ladder of [lang](../../../lang): a `lang.Parser` with the operators and the forms — if, fn, let, `~`, `()` — declared on it, and recursive descent around it for declarations and patterns; the top level binds in two passes, so the functions see each other — the mutual recursion of tml |
| `ast.go` | AST | s-expression printing, so tests read like trees; the nodes the lang ladder produces — literals, names, applications, infixes — are shared with lang as is |
| `types.go` | type syntax | mutable type variables, `prune`/`unify` with occurs check, `generalize`/`instantiate` |
| `infer.go` | AST → types | Hindley-Milner inference with let-polymorphism |
| `subst.go` | term machinery | the de Bruijn trio — `shift`, `subst`, `substitutionTops` — over terms and declarations |
| `value.go` | runtime values | the terms themselves, plus `VRec` — the recursive knot a declared function is |
| `eval.go` | AST → values | substitution interpreter: application is a literal beta step, and the vals of a scope leave the terms through substitution |
| `programs/` | the language itself | ten teaching programs, from literals to Church-encoded lists, each answering through its `main` |

The tests mirror the map (`lexer_test.go`, `parser_test.go`,
`infer_test.go`, `eval_test.go`, `programs_test.go`), and
`example_test.go` holds the pkg.go.dev examples.

## What you need to know: the type system

The typechecker is an implementation of the Hindley-Milner (HM) type
system — the same one behind Standard ML, Haskell and OCaml. Five
ideas carry all of it:

1. **Type variables and substitution.** A type is a graph of
   constructors (`int`, `bool`, arrows, tuples) and variables. A
   variable is a placeholder that can be *linked* to a type, and the
   link is mutable — assigning it is the substitution. In `types.go`
   the link lives in `TVar.Link`, and `prune` follows chains of links
   to their root, compressing them on the way.

2. **Unification.** To typecheck `if n = 0 then ...`, the checker must
   make two types equal: it *unifies* the type of `n` with `int`.
   `unify` in `types.go` links variables to types and recurses into
   arrows and tuples; a conflict (`int` vs `bool`) is a type error.

3. **The occurs check.** Unifying `'a` with `'a -> 'a` would build an
   infinite type — exactly what `fn x => x x` asks for. Before
   linking, `occurs` verifies the variable does not appear inside the
   other type, and rejects the program if it does.

4. **Generalization — where polymorphism comes from.** A `val`
   binding may be used at several types (`id 1` and `id true`):
   `generalize` in `types.go` quantifies every variable of the bound
   type that is not *locked* by the enclosing environment. Function
   parameters (`fn x => ...`) are never generalized — that is the
   whole difference between `val id = fn x => x` (polymorphic) and
   using a parameter (monomorphic).

5. **Instantiation.** At each use of a generalized binding,
   `instantiate` replaces the quantified variables with fresh ones,
   so the two uses of `id` constrain two independent copies.

One tml-specific extension: **top-level mutual recursion**. Classic
HM checks declarations in sequence; tml pre-binds every top-level
`fun` with a fresh shared variable before checking any body, so
`even` and `odd` can reference each other, and each still ends up
with the principal type `int -> bool` (`infer.go`, `Infer`).

The interpreter needs one new idea beyond substitution: the recursive
knot. Application is a literal beta step — `apply` matches the
argument pattern and moves the matched values into the body with
`substitutionTops` — and the `val`s of the top level and of every
`let` leave the terms the same way, so a running term carries no val
slots at all. What pure substitution cannot do is tie a function to
itself, so the functions of a scope — the top level included — become
one row of knots (`VRec` in `value.go`), each holding itself in its
own slot of the row; that single row is also what makes mutual
recursion work (`eval.go`, `evalProgram` and `apply`).

## Reading list

The primary sources, in order of how close they are to this code:

- **L. Cardelli, *Basic Polymorphic Typechecking* (1987)** — the
  mutation-based unification with path compression used here, down to
  the `prune` shape; the most implementation-honest of the classics.
- **L. Damas, R. Milner, *Principal Type-Schemes for Functional
  Programs* (1982)** — Algorithm W: the unification-based inference
  `infer.go` follows.
- **R. Milner, *A Theory of Type Polymorphism in Programming*
  (1978)** — the HM system itself: why these rules are sound and
  complete.
- **M. Grabmüller, *Algorithm W Step by Step* (2006)** — a gentle
  walkthrough of inference on paper, good before the papers.
- **S. Diehl, *Write You a Haskell*, the Hindley-Milner chapter** —
  inference to an AST very close to tml's, in a functional language.
- **B. C. Pierce, *Types and Programming Languages*, chapters 5–6** —
  nameless terms: the shift/subst machinery `subst.go` implements,
  down to the cutoff arithmetic.
- **H. Abelson, G. J. Sussman, *Structure and Interpretation of
  Computer Programs*, chapter 4** — the eval/apply interpreter
  pattern of `eval.go`.
- **R. Milner, M. Tofte, R. Harper, D. MacQueen, *The Definition of
  Standard ML* (Revised)** — where the surface syntax comes from.
- **D. Leijen, E. Meijer, *Parsec: Direct Style Parser Combinators
  for the Real World* (2001)** — the parsing style this whole module
  library builds on.

## Running it

The package is a library; the tests are the demo:

```bash
go test ./examples/lang/tml/ -run TestPrograms -v   # run the teaching programs
go test ./examples/lang/tml/ -run Example -v        # pkg.go.dev examples
```

Or from Go code:

```go
types, err := tml.Infer(src)          // map[string]Type
value, err := tml.Eval(src)           // the Value of main
value, err := tml.EvalFile(path)      // a program from a file
```

The teaching programs in `programs/` are the language tour — each one
explains a feature in comments, and `programs_test.go` checks what
every program computes.
