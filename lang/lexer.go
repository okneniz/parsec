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

// Definition describes the lexical syntax of a language. The type
// parameter S is the parser state the Custom combinators and the
// generated Lexer carry — the same state the expression machinery
// of a Parser carries, so a stepper sees what the parser sees.
type Definition[S parsec.Recoverable[S]] struct {
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
	Custom []func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, error)
}

func (def Definition[S]) withDefaults() Definition[S] {
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

func (def Definition[S]) keywordSet() map[string]struct{} {
	set := make(map[string]struct{}, len(def.Keywords))
	for _, kw := range def.Keywords {
		set[kw] = parsec.Stateless{}
	}
	return set
}

func (def Definition[S]) operatorTree() parsec.Tree[rune, strings.Position, Token, S] {
	cases := make(map[string]parsec.Combinator[rune, strings.Position, Token, S], len(def.Operators))

	for _, op := range def.Operators {
		cases[op] = parsec.Const[rune, strings.Position, Token, S](Token{
			Kind:   KindOperator,
			Lexeme: Lexeme(op),
		})
	}

	return parsec.NewLongestPrefixTree(cases, func(op string) []rune {
		return []rune(op)
	})
}

// Lexer turns the definition into a one-token stepper which carries
// the same state as the language parser.
func Lexer[S parsec.Recoverable[S]](def Definition[S]) tokens.Lexer[Kind, Lexeme, S] {
	def = def.withDefaults()
	token := def.token()
	trivia := def.trivia()

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, error) {
		if _, terr := trivia(state, buf); terr != nil {
			return Token{}, terr
		}

		if buf.IsEOF() {
			return Token{}, parsec.ErrEndOfFile
		}

		return nextToken(state, def, token, buf)
	}
}

// nextToken parses one token: custom combinators first, then the
// built-in dispatch by the leading rune.
func nextToken[S parsec.Recoverable[S]](
	state S,
	def Definition[S],
	token parsec.Combinator[rune, strings.Position, Token, S],
	buf parsec.Buffer[rune, strings.Position],
) (Token, error) {
	for _, custom := range def.Custom {
		pos := buf.Position()

		tok, err := custom(state, buf)
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

	tok, terr := token(state, buf)
	if terr != nil {
		return Token{}, terr
	}

	return tok, nil
}

func (def Definition[S]) token() parsec.Combinator[rune, strings.Position, Token, S] {
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

func (def Definition[S]) operator() parsec.Combinator[rune, strings.Position, Token, S] {
	operators := def.operatorTree()

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		op, lerr := operators.Lookup(state, buf)
		if lerr != nil {
			return Token{}, lerr
		}

		if op == nil {
			return Token{}, parsec.NewParseError(buf.Position(), "not an operator")
		}

		return op(state, buf)
	}
}

func (def Definition[S]) prefix() parsec.Combinator[rune, strings.Position, Token, S] {
	digits := strings.Many(0, strings.Try(
		strings.Satisfy[S]("expected digit", true, func(r rune) bool {
			return r >= '0' && r <= '9'
		}),
	))

	cases := map[rune]parsec.Combinator[rune, strings.Position, Token, S]{}

	if def.Integers {
		for _, d := range "0123456789" {
			digit := d

			cases[digit] = func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
				tail, _ := digits(state, buf)

				return Token{Kind: KindInt, Lexeme: Lexeme(string(digit) + string(tail))}, nil
			}
		}
	}

	for _, r := range def.Punctuation {
		punct := r

		cases[punct] = parsec.Const[rune, strings.Position, Token, S](Token{
			Kind:   KindSymbol,
			Lexeme: Lexeme(punct),
		})
	}

	return parsec.MapAs("no such leading rune", cases, strings.Any[S]())
}

func (def Definition[S]) stringForm(
	choice parsec.Combinator[rune, strings.Position, Token, S],
) parsec.Combinator[rune, strings.Position, Token, S] {
	quote := strings.Try(strings.Eq[S](`expected '"'`, '"'))
	str := stringRest[S](def.Escapes)

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		if _, qerr := quote(state, buf); qerr == nil {
			return str(state, buf)
		}

		return choice(state, buf)
	}
}

func (def Definition[S]) trivia() parsec.Combinator[rune, strings.Position, bool, S] {
	space := strings.Satisfy[S]("expected whitespace", true, def.Space)
	spaces := parsec.Many(0, strings.Try(space))

	var line parsec.Combinator[rune, strings.Position, bool, S]

	if def.LineComment != "" {
		line = tokens.LineComment[S](def.LineComment)
	}

	var blockOpen parsec.Combinator[rune, strings.Position, string, S]
	var blockBody parsec.Combinator[rune, strings.Position, bool, S]

	if def.BlockComment[0] != "" {
		blockOpen = strings.Try(strings.String[S]("comment start", def.BlockComment[0]))
		blockBody = tokens.CommentBody[S](def.BlockComment[0], def.BlockComment[1], def.NestedComments)
	}

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (bool, parsec.Error[strings.Position]) {
		for {
			if _, err := spaces(state, buf); err != nil {
				return false, err
			}

			// a line comment has no fatal errors of its own: the
			// only failure is the missing start marker
			if line != nil {
				if _, err := line(state, buf); err == nil {
					continue
				}
			}

			// a block comment commits after its opening
			if blockOpen != nil {
				if _, err := blockOpen(state, buf); err == nil {
					if _, cerr := blockBody(state, buf); cerr != nil {
						return false, cerr
					}

					continue
				}
			}

			return true, nil
		}
	}
}

func (def Definition[S]) identifier() parsec.Combinator[rune, strings.Position, Token, S] {
	first := strings.Satisfy[S]("expected identifier", true, def.IdentStart)
	rest := strings.Many(0, strings.Try(
		strings.Satisfy[S]("expected identifier", true, def.IdentLetter),
	))
	keywords := def.keywordSet()

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		head, err := first(state, buf)
		if err != nil {
			return Token{}, err
		}

		tail, _ := rest(state, buf)

		lexeme := Lexeme(string(head) + string(tail))

		kind := KindIdent
		if _, reserved := keywords[string(lexeme)]; reserved {
			kind = KindKeyword
		}

		return Token{Kind: kind, Lexeme: lexeme}, nil
	}
}

func stringRest[S any](escapes map[rune]rune) parsec.Combinator[rune, strings.Position, Token, S] {
	anyRune := strings.Any[S]()

	return func(state S, buf parsec.Buffer[rune, strings.Position]) (Token, parsec.Error[strings.Position]) {
		var text []rune

		for {
			if buf.IsEOF() {
				return Token{}, parsec.NewParseError(buf.Position(), "unterminated string literal")
			}

			r, rerr := anyRune(state, buf)
			if rerr != nil {
				return Token{}, rerr
			}

			if r == '"' {
				break
			}

			if r == '\\' {
				e, eerr := anyRune(state, buf)
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
