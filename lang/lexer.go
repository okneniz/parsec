package lang

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

type Token = tokens.Token[Kind, Lexeme]

type Position = strings.Position

type Kind string

type Lexeme string

const (
	KindKeyword  Kind = "keyword"
	KindIdent    Kind = "ident"
	KindOperator Kind = "operator"
	KindInt      Kind = "int"
	KindString   Kind = "string"
	KindSymbol   Kind = "symbol"
)

var ErrNoMatch = errors.New("no match")

type Definition struct {
	// LineComment starts a comment which runs to the end of the line,
	// as in "//" or "#"; empty disables line comments.
	LineComment string

	// BlockComment holds the start and end markers of a block comment,
	// as in "/*" "*/" or "(*" "*)"; an empty start disables block
	// comments. NestedComments makes blocks nest, SML and Pascal
	// style.
	BlockComment   [2]string
	NestedComments bool

	// Space decides which runes are whitespace trivia; the nil
	// default is unicode.IsSpace. A language where line breaks carry
	// structure narrows it, keeping the newline for its own token.
	Space func(rune) bool

	// IdentStart and IdentLetter decide which runes begin and continue
	// identifiers; nil defaults are a letter or underscore, then
	// letters, digits, underscores and apostrophes.
	IdentStart  func(rune) bool
	IdentLetter func(rune) bool

	// Keywords are reserved words, classified after lexing maximal
	// identifiers, so ifx stays an identifier next to if. Word
	// operators belong here too — operators must not start with an
	// identifier character.
	Keywords []string

	// Operators are matched by longest prefix: list == before = and
	// the lexer takes ==. A rune which starts some operator but
	// continues into none of them is a lexical error.
	Operators []string

	// Punctuation lists single-character tokens: brackets, commas,
	// semicolons. Each rune becomes one symbol token.
	Punctuation string

	// Integers enables decimal integer literals.
	Integers bool

	// Strings enables "..." literals with escapes, decoded into
	// Lexeme.
	Strings bool

	// Escapes maps the escape characters of string literals to the
	// runes they decode to; the nil default is \n, \t, \r, \\ and
	// \".
	Escapes map[rune]rune

	// Custom holds rune-level combinators for everything else, tried
	// in order before the built-in tokens. A combinator returns
	// ErrNoMatch to decline the input; any other error aborts Lex.
	// The tokens it produces may use arbitrary Kind strings.
	Custom []func(parsec.Buffer[rune, strings.Position]) (Token, error)
}

func (def Definition) withDefaults() Definition {
	if def.IdentStart == nil {
		def.IdentStart = func(r rune) bool {
			return unicode.IsLetter(r) || r == '_'
		}
	}

	if def.Space == nil {
		def.Space = unicode.IsSpace
	}

	if def.Escapes == nil {
		def.Escapes = map[rune]rune{
			'n':  '\n',
			't':  '\t',
			'r':  '\r',
			'\\': '\\',
			'"':  '"',
		}
	}

	if def.IdentLetter == nil {
		def.IdentLetter = func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '\''
		}
	}

	return def
}

func (def Definition) keywordSet() map[string]struct{} {
	set := make(map[string]struct{}, len(def.Keywords))
	for _, kw := range def.Keywords {
		set[kw] = struct{}{}
	}
	return set
}

func (def Definition) operatorTree() parsec.Tree[rune, strings.Position, Token] {
	cases := make(map[string]parsec.Combinator[rune, strings.Position, Token], len(def.Operators))

	for _, op := range def.Operators {
		cases[op] = parsec.Const[rune, strings.Position, Token](Token{
			Kind:   KindOperator,
			Lexeme: Lexeme(op),
		})
	}

	return parsec.NewLongestPrefixTree(cases, func(op string) []rune {
		return []rune(op)
	})
}

func Lexer(def Definition) tokens.Lexer[Kind, Lexeme] {
	def = def.withDefaults()
	token := def.token()
	trivia := def.trivia()

	return func(buf parsec.Buffer[rune, strings.Position]) (Token, error) {
		if _, terr := trivia(buf); terr != nil {
			return Token{}, terr
		}

		if buf.IsEOF() {
			return Token{}, parsec.ErrEndOfFile
		}

		return nextToken(def, token, buf)
	}
}

// nextToken parses one token: custom combinators first, then the
// built-in dispatch by the leading rune.
func nextToken(
	def Definition,
	token parsec.Combinator[rune, strings.Position, Token],
	buf parsec.Buffer[rune, strings.Position],
) (Token, error) {
	for _, custom := range def.Custom {
		pos := buf.Position()

		tok, err := custom(buf)
		if err == nil {
			return tok, nil
		}
		if !errors.Is(err, ErrNoMatch) {
			return Token{}, err
		}

		if err := buf.Seek(pos); err != nil {
			return Token{}, parsec.NewParseError(pos, err.Error())
		}
	}

	tok, terr := token(buf)
	if terr != nil {
		return Token{}, terr
	}

	return tok, nil
}

func (def Definition) token() parsec.Combinator[rune, strings.Position, Token] {
	choice := parsec.Choice("unrecognized character",
		parsec.Try(def.operator()),
		parsec.Try(def.prefix()),
		parsec.Try(def.identifier()),
	)

	if !def.Strings {
		return choice
	}

	return def.stringForm(choice)
}

func (def Definition) operator() parsec.Combinator[rune, strings.Position, Token] {
	operators := def.operatorTree()

	return func(buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		op, lerr := operators.Lookup(buf)
		if lerr != nil {
			return Token{}, lerr
		}

		if op == nil {
			return Token{}, parsec.NewParseError(buf.Position(), "not an operator")
		}

		return op(buf)
	}
}

func (def Definition) prefix() parsec.Combinator[rune, strings.Position, Token] {
	digits := strings.Many(0, strings.Try(
		strings.Satisfy("expected digit", true, func(r rune) bool {
			return r >= '0' && r <= '9'
		}),
	))

	cases := map[rune]parsec.Combinator[rune, strings.Position, Token]{}

	if def.Integers {
		for _, d := range "0123456789" {
			digit := d

			cases[digit] = func(buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
				tail, _ := digits(buf)

				return Token{Kind: KindInt, Lexeme: Lexeme(string(digit) + string(tail))}, nil
			}
		}
	}

	for _, r := range def.Punctuation {
		punct := r

		cases[punct] = parsec.Const[rune, strings.Position, Token](Token{
			Kind:   KindSymbol,
			Lexeme: Lexeme(punct),
		})
	}

	return parsec.MapAs("no such leading rune", cases, strings.Any())
}

func (def Definition) stringForm(
	choice parsec.Combinator[rune, strings.Position, Token],
) parsec.Combinator[rune, strings.Position, Token] {
	quote := strings.Try(strings.Eq(`expected '"'`, '"'))
	str := stringRest(def.Escapes)

	return func(buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		if _, qerr := quote(buf); qerr == nil {
			return str(buf)
		}

		return choice(buf)
	}
}

func (def Definition) trivia() parsec.Combinator[rune, strings.Position, bool] {
	space := strings.Satisfy("expected whitespace", true, def.Space)
	spaces := parsec.Many(0, strings.Try(space))

	var line parsec.Combinator[rune, strings.Position, bool]

	if def.LineComment != "" {
		line = tokens.LineComment(def.LineComment)
	}

	var blockOpen parsec.Combinator[rune, strings.Position, string]
	var blockBody parsec.Combinator[rune, strings.Position, bool]

	if def.BlockComment[0] != "" {
		blockOpen = strings.Try(strings.String("comment start", def.BlockComment[0]))
		blockBody = tokens.CommentBody(def.BlockComment[0], def.BlockComment[1], def.NestedComments)
	}

	return func(buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		for {
			if _, err := spaces(buf); err != nil {
				return false, err
			}

			// a line comment has no fatal errors of its own: the
			// only failure is the missing start marker
			if line != nil {
				if _, err := line(buf); err == nil {
					continue
				}
			}

			// a block comment commits after its opening
			if blockOpen != nil {
				if _, err := blockOpen(buf); err == nil {
					if _, cerr := blockBody(buf); cerr != nil {
						return false, cerr
					}

					continue
				}
			}

			return true, nil
		}
	}
}

func (def Definition) identifier() parsec.Combinator[rune, strings.Position, Token] {
	first := strings.Satisfy("expected identifier", true, def.IdentStart)
	rest := strings.Many(0, strings.Try(
		strings.Satisfy("expected identifier", true, def.IdentLetter),
	))
	keywords := def.keywordSet()

	return func(buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		head, err := first(buf)
		if err != nil {
			return Token{}, err
		}

		tail, _ := rest(buf)

		lexeme := Lexeme(string(head) + string(tail))

		kind := KindIdent
		if _, reserved := keywords[string(lexeme)]; reserved {
			kind = KindKeyword
		}

		return Token{Kind: kind, Lexeme: lexeme}, nil
	}
}

func stringRest(escapes map[rune]rune) parsec.Combinator[rune, strings.Position, Token] {
	any := strings.Any()

	return func(buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		var text []rune

		for {
			if buf.IsEOF() {
				return Token{}, parsec.NewParseError(buf.Position(), "unterminated string literal")
			}

			r, rerr := any(buf)
			if rerr != nil {
				return Token{}, rerr
			}

			if r == '"' {
				break
			}

			if r == '\\' {
				e, eerr := any(buf)
				if eerr != nil {
					return Token{}, eerr
				}

				decoded, known := escapes[e]
				if !known {
					return Token{}, parsec.NewParseError(buf.Position(), fmt.Sprintf("unknown escape sequence \\%c", e))
				}
				r = decoded
			}

			text = append(text, r)
		}

		return Token{
			Kind:   KindString,
			Lexeme: Lexeme(text),
		}, nil
	}
}
