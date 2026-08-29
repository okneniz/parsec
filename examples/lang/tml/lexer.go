package tml

import (
	"errors"
	"unicode"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/lang"
	"github.com/okneniz/parsec/strings"
)

// Token, Kind and Lexeme are the token vocabulary of the lang
// package: the mapping from runes to tokens is configured there, in
// one place, and this package only consumes it.
type (
	Token    = lang.Token
	Kind     = lang.Kind
	Lexeme   = lang.Lexeme
	Position = lang.Position
)

// The token kinds the generated lexer produces for tiny ml.
const (
	KindKeyword  = lang.KindKeyword
	KindIdent    = lang.KindIdent
	KindOperator = lang.KindOperator
	KindInt      = lang.KindInt
	KindSymbol   = lang.KindSymbol
)

// symbolicChars is the set of characters symbolic identifiers are
// built from, as in Standard ML: every operator is a run of these.
const symbolicChars = "!%&$#+-/:<=>?@\\~`^|*"

var symbolicRun = strings.Some(0, "expected symbolic identifier",
	strings.Try(strings.OneOf("expected symbolic character", []rune(symbolicChars)...)),
)

// symbolicIdentifier parses the maximal run of symbolic characters,
// as in Standard ML, and returns it as an operator token waiting for
// the fixity table. The reserved = and => are declined with
// lang.ErrNoMatch: the definition below lists them as its operators,
// and the lexer tries Custom before its own branches.
func symbolicIdentifier(buf parsec.Buffer[rune, strings.Position]) (Token, error) {
	runes, err := symbolicRun(buf)
	if err != nil {
		return Token{}, lang.ErrNoMatch
	}

	lexeme := Lexeme(runes)

	if lexeme == "=" || lexeme == "=>" {
		return Token{}, lang.ErrNoMatch
	}

	return Token{Kind: lang.KindOperator, Lexeme: lexeme}, nil
}

// definition describes the lexical structure of tiny ml: the nested
// comment, the reserved words, the integer literals, the two reserved
// symbols and the punctuation. Everything a Definition cannot
// describe — the open world of symbolic operators — goes into
// Custom.
var definition = lang.Definition{
	BlockComment:   [2]string{"(*", "*)"},
	NestedComments: true,
	IdentStart:     unicode.IsLetter,
	Keywords: []string{
		"andalso", "else", "end", "false", "fn", "fun",
		"if", "in", "let", "orelse", "then", "true", "val",
	},
	Operators:   []string{"=>", "="},
	Punctuation: "_(),",
	Integers:    true,
	Custom:      []func(parsec.Buffer[rune, strings.Position]) (Token, error){symbolicIdentifier},
}

// lexer is the one-token stepper of tiny ml, generated from the
// definition by the lang package.
var lexer = lang.Lexer(definition)

// Lex splits src into tiny ml tokens by running the generated lexer
// to the end of input. Whitespace and comments are skipped.
func Lex(src string) ([]Token, error) {
	buf := strings.Buffer([]rune(src))
	stream := make([]Token, 0, 64)

	for {
		tok, err := lexer(buf)
		if errors.Is(err, parsec.ErrEndOfFile) {
			break
		}
		if err != nil {
			return nil, err
		}

		stream = append(stream, tok)
	}

	return stream, nil
}
