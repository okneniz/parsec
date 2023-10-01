package strings

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/okneniz/parsec/testing"
)

func TestEq(t *testing.T) {
	comb := Eq('c')

	result, err := ParseString("a", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("b", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')
}

func TestNotEq(t *testing.T) {
	comb := NotEq('c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("abc", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestOneOf(t *testing.T) {
	comb := OneOf('a', 'b', 'c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')

	result, err = ParseString("d", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestSequenceOf(t *testing.T) {
	comb := SequenceOf('f', 'o', 'o')

	result, err := ParseString("foo", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{'f', 'o', 'o'})

	result, err = ParseString("foobar", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{'f', 'o', 'o'})

	result, err = ParseString("fo", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString(" foobar", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString(" ", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}

func TestMap(t *testing.T) {
	cases := map[rune]int{'a': 1, 'b': 2, 'c': 3}

	comb := Some(
		1,
		SkipMany(
			NoneOf('a', 'b', 'c'),
			Map(cases, Any()),
		),
	)

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1})

	result, err = ParseString("..a//b++c**d,,e--a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1, 2, 3, 1})

	result, err = ParseString("bb", comb)
	Check(t, err)
	AssertSlice(t, result, []int{2, 2})

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}

func TestMapStrings(t *testing.T) {
	cases := map[string]int{"a": 1, "b": 2, "c": 3}

	comb := Some(
		1,
		SkipMany(
			NoneOf('a', 'b', 'c'),
			MapStrings(cases),
		),
	)

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1})

	result, err = ParseString("..a//b++c**d,,e--a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1, 2, 3, 1})

	result, err = ParseString("bb", comb)
	Check(t, err)
	AssertSlice(t, result, []int{2, 2})

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := String("foo")

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("foobar", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("bar", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("baz", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString(" foo", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString(" foobar", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestOneOfStrings(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := OneOfStrings("foo", "bar")

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("foobar", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("barbaz", comb)
		Check(t, err)
		AssertEq(t, result, "bar")

		result, err = ParseString("baz", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Many(
			0,
			SkipMany(
				NoneOf('f', 'o', 'b', 'a', 'r'),
				OneOfStrings("foo", "bar"),
			),
		)

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"foo"})

		result, err = ParseString("barfoo", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("bar12334foo123", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("bar12334foo123baz", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("12311231820398", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func BenchmarkMap(b *testing.B) {
	seed := time.Now().UnixNano()
	source := rand.New(rand.NewSource(seed))
	r := rand.New(source)

	b.Log("seed: ", seed)

	dict := map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"May": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Aug": time.August,
		"Sep": time.September,
		"Oct": time.October,
		"Nov": time.November,
		"Dec": time.December,
	}

	gen := func(count int) []string {
		examples := make([]string, 0, count)

		for {
			for key := range dict {
				examples = append(examples, key)
				count--

				if count == 0 {
					break
				}
			}

			if count == 0 {
				break
			}
		}

		r.Shuffle(len(examples), func(i, j int) { examples[i], examples[j] = examples[j], examples[i] })

		return examples
	}

	b.Run("MapString", func(b *testing.B) {
		examples := gen(b.N)
		comb := MapStrings(dict)

		b.ResetTimer()

		for _, example := range examples {
			_, _ = ParseString(example, comb)
		}
	})

	b.Run("Map", func(b *testing.B) {
		examples := gen(b.N)
		comb := Map(dict, Cast(
			Count(3, Any()),
			func(x []rune) (string, error) {
				return string(x), nil
			},
		))

		b.ResetTimer()

		for _, example := range examples {
			_, _ = ParseString(example, comb)
		}
	})
}
