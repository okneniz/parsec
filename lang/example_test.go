// This file holds the runnable documentation examples: they appear
// on pkg.go.dev as the package examples.
package lang_test

import (
	"fmt"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/lang"
	"github.com/okneniz/parsec/strings"
	"github.com/okneniz/parsec/tokens"
)

func Example() {
	def := lang.Definition{
		LineComment:    "//",
		BlockComment:   [2]string{"(*", "*)"},
		NestedComments: true,

		Keywords:    []string{"let", "in"},
		Operators:   []string{"+", "-", "*", "/"},
		Punctuation: "()=",
		Integers:    true,
	}

	p := lang.NewParser().
		Infix(1, "+", "-").
		Infix(2, "*", "/").
		Prefix("-")

	expression := "1 + 2 * -3 (* nested (* comment *) *)"
	buf := strings.Buffer([]rune(expression))
	ml := lang.New(def, p)

	e, perr := ml(buf)
	if perr != nil {
		panic(perr)
	}

	fmt.Println(e)
	// Output: (+ 1 (* 2 (- 3)))
}

// A sizeof form: a keyword-led operand that binds like a prefix
// operator, tighter than any infix one.
func ExampleParser_PrefixForm() {
	def := lang.Definition{
		Keywords:  []string{"sizeof"},
		Operators: []string{"+", "*"},
		Integers:  true,
	}

	sizeof := func(l lang.Ladder) parsec.Combinator[rune, strings.Position, lang.Expr] {
		kw := tokens.Exact(l.Lex, lang.KindKeyword, "sizeof")

		return func(buf parsec.Buffer[rune, strings.Position]) (lang.Expr, parsec.Error[strings.Position]) {
			if _, err := kw(buf); err != nil {
				return nil, err
			}

			sub, serr := l.Atom(buf)
			if serr != nil {
				return nil, serr
			}

			return lang.Unary{Op: "sizeof", Sub: sub}, nil
		}
	}

	p := lang.NewParser().
		Infix(1, "+").
		Infix(2, "*").
		PrefixForm(sizeof)

	e, perr := lang.New(def, p)(strings.Buffer([]rune("1 + sizeof x * 2")))
	if perr != nil {
		panic(perr)
	}

	fmt.Println(e)
	// Output: (+ 1 (* (sizeof x) 2))
}
