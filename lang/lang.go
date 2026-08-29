package lang

import (
	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// New returns a combinator which parses one complete expression from
// a rune buffer with p: an operand chain, up to the end of input.
// Tokens are pulled from the input with the Lexer of def — the two
// stages share one buffer, and the rune position is the token
// position. The extras are the extra operands of Parser.Expr:
//
//	expr, err := lang.New(def, p, let)(
//		strings.Buffer([]rune(src)),
//	)
func New(
	def Definition,
	p *Parser,
	extras ...Operand,
) parsec.Combinator[rune, strings.Position, Expr] {
	def = def.withDefaults()
	lex := Lexer(def)
	expr := p.Expr(lex, extras...)
	unexpected := parsec.Unexpected(lex, "expected end of input")
	trivia := def.trivia()

	return func(buf parsec.Buffer[rune, strings.Position]) (Expr, parsec.Error[strings.Position]) {
		e, err := expr(buf)
		if err != nil {
			return nil, located(buf, err)
		}

		// the end check asks the lexer's own trivia skipper: the
		// trailing whitespace and comments are not input
		if _, terr := trivia(buf); terr != nil {
			return nil, located(buf, terr)
		}

		if !buf.IsEOF() {
			_, uerr := unexpected(buf)

			return nil, located(buf, uerr)
		}

		return e, nil
	}
}

// located converts an error to the parsec.Error the combinator must
// return; the errors raised by the parser already carry positions.
func located(buf parsec.Buffer[rune, strings.Position], err error) parsec.Error[strings.Position] {
	if parseErr, ok := err.(parsec.Error[strings.Position]); ok {
		return parseErr
	}

	return parsec.NewParseError(buf.Position(), err.Error())
}
