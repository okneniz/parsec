package parsec

import (
	"golang.org/x/exp/constraints"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestCards(t *testing.T) {
	digit := Range[byte, int]('0', '9')
	notDigit := NotRange[byte, int]('0', '9')

	// from https://www.regular-expressions.info/creditcard.html

	visa := Trace(t, "visa",
		Cast(
			Sequence(
				16,
				Count(1, Eq[byte, int]('4')),
				Count(12, digit),
				Optional(Count(3, digit), []byte{}),
			),
			toString,
		),
	)

	master := Trace(t, "master",
		Cast(
			Sequence(
				15,
				Choice(
					Trace(t, "master 1",
						Concat(
							4,
							Count(1, Eq[byte, int]('5')),
							Count(1, Range[byte, int]('1', '5')),
							Count(2, Range[byte, int]('0', '9')),
						),
					),
					Trace(t, "master 2",
						Concat(
							4,
							Count(3, Eq[byte, int]('2')),
							Count(1, Range[byte, int]('1', '9')),
						),
					),
					Trace(t, "master 3",
						Concat(
							4,
							Count(2, Eq[byte, int]('2')),
							Count(1, Range[byte, int]('3', '9')),
							Count(1, Range[byte, int]('0', '9')),
						),
					),
					Trace(t, "master 4",
						Concat(
							4,
							Count(1, Eq[byte, int]('2')),
							Count(1, Range[byte, int]('3', '6')),
							Count(2, Range[byte, int]('0', '9')),
						),
					),
					Trace(t, "master 5",
						Concat(
							4,
							Count(1, Eq[byte, int]('2')),
							Count(1, Range[byte, int]('3', '6')),
							Count(2, Range[byte, int]('0', '9')),
						),
					),
					Trace(t, "master 5",
						Sequence(
							4,
							Eq[byte, int]('2'),
							Eq[byte, int]('7'),
							OneOf[byte, int]('0', '1'),
							Range[byte, int]('0', '9'),
						),
					),
					Trace(t, "master 6",
						SequenceOf[byte, int]('2', '7', '2', '0'),
					),
				),
				Count(12, Range[byte, int]('0', '9')),
			),
			toString,
		),
	)

	americanExpress := Trace(t, "american express",
		Cast(
			Sequence(
				2,
				Sequence(
					2,
					Eq[byte, int]('3'),
					OneOf[byte, int]('4', '7'),
				),
				Count(13, Range[byte, int]('0', '9')),
			),
			toString,
		),
	)

	cards := Choice(
		Try(visa),
		Try(master),
		Try(americanExpress),
	)

	noice := Many(10, Try(notDigit))
	comb := Skip(noice, SepEndBy(4, cards, noice))

	cardNums := []string{
		"4111111111111111",
		"4012888888881881",
		"4222222222222",
		"5555555555554444",
		"5105105105105100",
		"378282246310005",
		"371449635398431",
		"378734493671000",
	}

	seed := time.Now().UnixNano()
	t.Log("seed: ", seed)
	shuffle := shuffler[string](seed)

	input := joinBy(
		noicer(seed, byte('0'), byte('9')),
		shuffle(cardNums)...,
	)

	t.Log("input:")
	t.Logf("%#v", input)

	result, err := ParseBytes([]byte(input), comb)
	check(t, err)
	assertSlice(t, sorted(result...), sorted(cardNums...))
}

func copyOf[T any](data []T) []T {
	result := make([]T, len(data))
	for i, x := range data {
		result[i] = x
	}
	return result
}

func shuffler[T any](seed int64) func([]T) []T {
	source := rand.New(rand.NewSource(seed))

	return func(data []T) []T {
		result := copyOf(data)

		source.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})

		return result
	}
}

func sorted[T constraints.Ordered](data ...T) []T {
	result := copyOf(data)

	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

func joinBy(f func() string, data ...string) string {
	result := ""

	for i := 0; i < len(data)-1; i++ {
		result += data[i]
		result += f()
	}

	result += data[len(data)-1]
	return result
}

func noicer(seed int64, from, to byte) func() string {
	source := rand.New(rand.NewSource(seed))

	return func() string {
		size := 1 + source.Intn(9)
		result := ""

		for i := 0; i < size; i++ {
			var b byte

			if source.Intn(10)%2 == 0 {
				b = byte(1 + source.Intn(int(from)-1))
			} else {
				b = byte(int(to) + 1 + source.Intn(math.MaxUint8-int(to)))
			}

			result += string(b)
		}

		return result
	}
}

func toString(xs [][]byte) (string, error) {
	s := ""
	for _, x := range xs {
		s += string(x)
	}
	return s, nil
}
