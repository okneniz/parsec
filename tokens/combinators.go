package tokens

import (
	"fmt"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// Lexer is a one-token stepper over a rune buffer: it reads the next
// token, skipping the leading whitespace and comments itself, and
// reports the exhausted input with parsec.ErrEndOfFile — the same
// sentinel the buffers return.
type Lexer[K comparable, L comparable] func(buf parsec.Buffer[rune, strings.Position]) (Token[K, L], error)

var anyRune = strings.Any()

// LineComment parses one line comment: the start marker, then every
// rune up to and including the end of the line. A line comment has
// no fatal errors of its own — the only failure is the missing start
// marker, with the input left intact.
func LineComment(start string) parsec.Combinator[rune, strings.Position, bool] {
	mark := strings.Try(strings.String("line comment", start))

	return func(buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		if _, err := mark(buf); err != nil {
			return false, err
		}

		for !buf.IsEOF() {
			r, err := anyRune(buf)
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
func CommentBody(
	open string,
	end string,
	nested bool,
) parsec.Combinator[rune, strings.Position, bool] {
	openMark := strings.Try(strings.String("comment start", open))
	endMark := strings.Try(strings.String("comment end", end))

	return func(buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		for depth := 1; depth > 0; {
			if buf.IsEOF() {
				return false, parsec.NewParseError(buf.Position(), "unterminated comment")
			}

			if _, err := endMark(buf); err == nil {
				depth--
				continue
			}

			if nested {
				if _, err := openMark(buf); err == nil {
					depth++
					continue
				}
			}

			if _, err := anyRune(buf); err != nil {
				return false, err
			}
		}

		return true, nil
	}
}

// Satisfy matches the next token by predicate. It consumes the token
// through lex even when the predicate fails — the greedy semantics
// of parsec.Satisfy — so wrap speculative uses in parsec.Try.
func Satisfy[K comparable, L comparable](
	lex Lexer[K, L],
	msg string,
	f func(Token[K, L]) bool,
) parsec.Combinator[rune, strings.Position, Token[K, L]] {
	return func(buf parsec.Buffer[rune, strings.Position]) (Token[K, L], parsec.Error[strings.Position]) {
		pos := buf.Position()

		tok, err := lex(buf)
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
func Exact[K comparable, L comparable](
	lex Lexer[K, L],
	kind K,
	lexeme L,
) parsec.Combinator[rune, strings.Position, Token[K, L]] {
	return Satisfy(
		lex,
		fmt.Sprintf("expected %v %q", kind, fmt.Sprint(lexeme)),
		func(t Token[K, L]) bool {
			return t.Kind == kind && t.Lexeme == lexeme
		},
	)
}

// OfKind matches the next token of the given kind.
func OfKind[K comparable, L comparable](
	lex Lexer[K, L],
	kind K,
) parsec.Combinator[rune, strings.Position, Token[K, L]] {
	return Satisfy(
		lex,
		fmt.Sprintf("expected %v", kind),
		func(t Token[K, L]) bool {
			return t.Kind == kind
		},
	)
}
