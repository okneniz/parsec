// Package lang assembles parsers for basic programming languages from
// a declarative description, in the spirit of the LanguageDef record
// and makeTokenParser of Haskell parsec — adapted to the two-stage
// pipeline of this library: a rune-level lexer generated from a
// Definition produces tokens over the same rune buffer, and the
// expression machinery reads them through it.
//
// A Definition describes what simple languages share — which runes
// are whitespace, line and nested block comments, identifiers with
// reserved words, operators with longest match, integer and string
// literals with their escapes, punctuation — and Lexer turns it into
// a one-token stepper. The stepper dispatches by prefix, with no
// lookahead: the operators through a longest-prefix tree, the
// leading rune of everything else through a table, and the
// identifier run as the open-domain fallback; a string literal is
// guarded ahead of the dispatch, so that its body errors stay fatal.
// Custom holds rune-level combinators for the forms a Definition
// cannot describe — char literals, hexadecimal numbers, SML-style
// symbolic operators, the structural tokens of an indentation
// language — a custom declines input by returning ErrNoMatch.
//
// New assembles the full parser from a Definition and a Parser: a
// combinator over the rune buffer, parsec.Combinator[rune,
// strings.Position, Expr], composable with the rest of the library.
// The Parser holds a small fixity table; its Expr parses prefix,
// postfix and infix operator chains with the declared precedences
// and associativities. An operator is a lexeme with fixity status,
// whatever kind the lexer gave its token — the reserved words
// andalso and orelse of sml are infix operators of their own level.
//
// The forms a language has beyond operators — if e then e else e
// above the infix chain, a sizeof e or a let among the operands —
// rarely reduce to an operator lexeme. ExprForm declares the ones of
// the expression level, PrefixForm the ones of the operand level,
// and the Ladder each form receives hands it the recursion: Expr is
// the full expression, forms included, Atom one operand, Lex the
// lexer the ladder reads through. The extras of Parser.Expr remain
// the plain token-level operand combinators for the rest of the
// syntax: boolean literals, calls, indexing.
package lang
