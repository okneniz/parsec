package tokens

import (
	"fmt"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// Lexer is a one-token stepper over a rune buffer: it reads the next
// token, skipping the leading whitespace and comments itself, and
// reports the exhausted input with parsec.ErrEndOfFile — the same
// sentinel the buffers return. The stepper carries the parser state
// of type S, the same value the combinators carry.
type Lexer[K comparable, L comparable, S any] func(state S, buf parsec.Buffer[rune, strings.Position]) (Token[K, L], error)

var anyRune = strings.Any[parsec.Stateless]()

// LineComment parses one line comment: the start marker, then every
// rune up to and including the end of the line. A line comment has
// no fatal errors of its own — the only failure is the missing start
// marker, with the input left intact.
func LineComment[S any](start string) parsec.Combinator[rune, strings.Position, bool, S] {
	mark := strings.Try(strings.String[S]("line comment", start))

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		if _, err := mark(state, buf); err != nil {
			return false, err
		}

		for !buf.IsEOF() {
			r, err := anyRune(parsec.Stateless{}, buf)
			if err != nil {
				return false, err
			}

			if r == '\n' {
				return true, nil
			}
		}

		return true, nil
	}
}

// CommentBody parses the rest of a block comment after its opening
// was consumed by the caller: everything up to the matching close,
// counting the nesting depth when nested — a manual loop over the
// buffer, because nesting needs a counter which the stock
// combinators do not provide. An unterminated comment fails; the
// caller keeps the opening probe separate, so a missing comment
// stays distinguishable from a broken one.
func CommentBody[S any](
	open string,
	end string,
	nested bool,
) parsec.Combinator[rune, strings.Position, bool, S] {
	openMark := strings.Try(strings.String[S]("comment start", open))
	endMark := strings.Try(strings.String[S]("comment end", end))

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		for depth := 1; depth > 0; {
			if buf.IsEOF() {
				return false, parsec.NewParseError(buf.Position(), "unterminated comment")
			}

			if _, err := endMark(state, buf); err == nil {
				depth--
				continue
			}

			if nested {
				if _, err := openMark(state, buf); err == nil {
					depth++
					continue
				}
			}

			if _, err := anyRune(parsec.Stateless{}, buf); err != nil {
				return false, err
			}
		}

		return true, nil
	}
}

// Satisfy matches the next token by predicate. It consumes the token
// through lex even when the predicate fails — the greedy semantics
// of parsec.Satisfy — so wrap speculative uses in parsec.Try.
func Satisfy[K comparable, L comparable, S any](
	lex Lexer[K, L, S],
	msg string,
	f func(Token[K, L]) bool,
) parsec.Combinator[rune, strings.Position, Token[K, L], S] {
	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token[K, L], parsec.Error[strings.Position]) {
		pos := buf.Position()

		tok, err := lex(state, buf)
		if err != nil {
			// a positional lexical error is chained, not flattened:
			// its type and message stay reachable through Previous
			if parseErr, ok := err.(parsec.Error[strings.Position]); ok {
				return tok, parsec.NewParseError(pos, msg, parseErr)
			}

			return tok, parsec.NewParseError(pos, msg)
		}

		if f(tok) {
			return tok, nil
		}

		return tok, parsec.NewParseError(pos, msg)
	}
}

// Exact matches a token of the given kind with the given lexeme.
func Exact[K comparable, L comparable, S any](
	lex Lexer[K, L, S],
	kind K,
	lexeme L,
) parsec.Combinator[rune, strings.Position, Token[K, L], S] {
	return Satisfy(
		lex,
		fmt.Sprintf("expected %v %q", kind, fmt.Sprint(lexeme)),
		func(t Token[K, L]) bool {
			return t.Kind == kind && t.Lexeme == lexeme
		},
	)
}

// OfKind matches the next token of the given kind.
func OfKind[K comparable, L comparable, S any](
	lex Lexer[K, L, S],
	kind K,
) parsec.Combinator[rune, strings.Position, Token[K, L], S] {
	return Satisfy(
		lex,
		fmt.Sprintf("expected %v", kind),
		func(t Token[K, L]) bool {
			return t.Kind == kind
		},
	)
}
