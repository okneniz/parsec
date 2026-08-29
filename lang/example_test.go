// This file holds the runnable documentation examples: they appear
// on pkg.go.dev as the package examples.
package lang_test

import (
	"fmt"

	"github.com/okneniz/parsec/lang"
	"github.com/okneniz/parsec/strings"
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
