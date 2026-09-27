package strings_test

import (
	"fmt"

	"github.com/okneniz/parsec"
	. "github.com/okneniz/parsec/strings"
)

// The simplest parser: an unsigned integer with optional
// surrounding whitespaces.
func Example() {
	parser := Padded(
		Try(Space[parsec.Stateless]("whitespace")),
		Unsigned[int, parsec.Stateless](),
	)

	result, err := ParseString(parsec.Stateless{}, " 42 ", parser)
	fmt.Println(result, err)
	// Output: 42 <nil>
}

// Choice tries alternatives one by one; each alternative
// must be wrapped in Try to rewind the buffer on failure.
func ExampleChoice() {
	keyword := Choice(
		"expected keyword",
		Try(String[parsec.Stateless]("expected 'true'", "true")),
		Try(String[parsec.Stateless]("expected 'false'", "false")),
	)

	result, err := ParseString(parsec.Stateless{}, "false", keyword)
	fmt.Println(result, err)
	// Output: false <nil>
}

// Satisfy matches a single rune against an arbitrary predicate.
func ExampleSatisfy() {
	vowel := Satisfy[parsec.Stateless]("expected vowel", true, func(r rune) bool {
		return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
	})

	result, err := ParseString(parsec.Stateless{}, "u", vowel)
	fmt.Printf("%c %v\n", result, err)
	// Output: u <nil>
}

// Range matches a single rune inside an inclusive range,
// for example a decimal digit.
func ExampleRange() {
	digit := Range[parsec.Stateless]("expected decimal digit", '0', '9')

	result, err := ParseString(parsec.Stateless{}, "7", digit)
	fmt.Printf("%c %v\n", result, err)
	// Output: 7 <nil>
}

// SequenceOf expects a fixed sequence of runes, like a keyword.
func ExampleSequenceOf() {
	let := SequenceOf[parsec.Stateless]("expected 'let'", []rune("let")...)

	result, err := ParseString(parsec.Stateless{}, "let x = 1", let)
	fmt.Println(string(result), err)
	// Output: let <nil>
}

// Many applies a combinator as many times as possible.
func ExampleMany() {
	digits := Many(4, Digit[parsec.Stateless]("expected digit"))

	result, err := ParseString(parsec.Stateless{}, "0123", digits)
	fmt.Println(string(result), err)
	// Output: 0123 <nil>
}

// Optional falls back to a default value instead of failing.
func ExampleOptional() {
	signed := And(
		Optional(Try(Eq[parsec.Stateless]("expected sign", '-')), '+'),
		Unsigned[int, parsec.Stateless](),
		func(sign rune, num int) int {
			if sign == '-' {
				return -num
			}
			return num
		},
	)

	for _, input := range []string{"-42", "42"} {
		result, err := ParseString(parsec.Stateless{}, input, signed)
		fmt.Println(result, err)
	}

	// Output:
	// -42 <nil>
	// 42 <nil>
}

// SepBy parses a separator-separated list.
func ExampleSepBy() {
	list := SepBy(3, Unsigned[int, parsec.Stateless](), Comma[parsec.Stateless]())

	result, err := ParseString(parsec.Stateless{}, "1,2,3", list)
	fmt.Println(result, err)
	// Output: [1 2 3] <nil>
}

// Parens parses the body between '(' and ')'.
// The body element is wrapped in Try: a greedy combinator consumes
// the rune it failed on, which would otherwise eat the closing paren.
func ExampleParens() {
	ident := Some(4, "expected identifier", Try(Letter[parsec.Stateless]("expected letter")))

	result, err := ParseString(parsec.Stateless{}, "(abc)", Parens(ident))
	fmt.Println(string(result), err)
	// Output: abc <nil>
}

// MapStrings translates the parsed text using a map,
// matching the longest key on a prefix collision.
func ExampleMapStrings() {
	literal := MapStrings[bool, parsec.Stateless]("expected boolean literal", map[string]bool{
		"true":  true,
		"false": false,
	})

	result, err := ParseString(parsec.Stateless{}, "true", literal)
	fmt.Println(result, err)
	// Output: true <nil>
}

// Chainl1 parses left-associative binary expressions:
// 10 - 4 - 2 is (10 - 4) - 2, not 10 - (4 - 2).
func ExampleChainl1() {
	number := Padded(Try(Space[parsec.Stateless]("whitespace")), Unsigned[int, parsec.Stateless]())

	minus := Padded(Try(Space[parsec.Stateless]("whitespace")), Eq[parsec.Stateless]("expected '-'", '-'))
	op := Cast(minus, func(r rune) (parsec.BinaryOp[int], error) {
		return func(a, b int) int { return a - b }, nil
	})

	result, err := ParseString(parsec.Stateless{}, "10 - 4 - 2", Chainl1(number, op))
	fmt.Println(result, err)
	// Output: 4 <nil>
}

// Seq runs a combinator as a lazy sequence over the buffer: every
// successful step yields its value, the first failure yields its
// error and ends the sequence, and the buffer stays in place for the
// next combinator.
func ExampleSeq() {
	buf := Buffer([]rune("aab"))

	seq, _ := Seq(Try(Eq[parsec.Stateless]("expected a", 'a')))(parsec.Stateless{}, buf)

	for r, err := range seq {
		fmt.Println(r, err)
	}

	next, err := Try(Eq[parsec.Stateless]("expected b", 'b'))(parsec.Stateless{}, buf)
	fmt.Println(next, err)

	// Output:
	// 97 <nil>
	// 97 <nil>
	// 0 Parse error at line=0 column=2 index=2: expected a
	// 98 <nil>
}
