package testing

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"golang.org/x/exp/constraints"
)

func Check(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	} else {
		t.Log("catch error: ", err)
	}
}

func Assert(t *testing.T, x bool, m string) {
	t.Helper()

	if !x {
		t.Fatal(m)
	}
}

func AssertEq[T comparable](t *testing.T, x, y T) {
	t.Helper()

	if x != y {
		t.Fatalf("%v != %v", x, y)
	}
}

func AssertSlice[T comparable](t *testing.T, xs, ys []T) {
	t.Helper()

	if len(xs) != len(ys) {
		t.Fatalf("%v != %v", xs, ys)
	}

	for i, x := range xs {
		if x != ys[i] {
			t.Fatalf("%v != %v", xs, ys)
		}
	}
}

func AssertEqDump[T any](t *testing.T, actual, expected T) {
	t.Helper()

	ex, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	ac, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if string(ex) != string(ac) {
		t.Errorf("expected %v", string(ex))
		t.Errorf("actual %v", string(ac))
		t.Fatal("invalid result")
	}
}

func CopyOf[T any](data []T) []T {
	result := make([]T, len(data))
	copy(result, data)
	return result
}

func Shuffler[T any](seed int64) func([]T) []T {
	source := rand.New(rand.NewSource(seed))

	return func(data []T) []T {
		result := CopyOf(data)

		source.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})

		return result
	}
}

func Sorted[T constraints.Ordered](data ...T) []T {
	result := CopyOf(data)

	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

func JoinBy(f func() string, data ...string) string {
	result := ""

	for i := 0; i < len(data)-1; i++ {
		result += data[i]
		result += f()
	}

	result += data[len(data)-1]
	return result
}

func Noicer(seed int64, from, to rune) func() string {
	source := rand.New(rand.NewSource(seed))

	return func() string {
		size := 1 + source.Intn(9)
		result := ""

		for i := 0; i < size; i++ {
			var b rune

			if source.Intn(10)%2 == 0 {
				b = rune(1 + source.Intn(int(from)-1))
			} else {
				b = rune(int(to) + 1 + source.Intn(math.MaxUint8-int(to)))
			}

			result += string(b)
		}

		return result
	}
}

func DigitsToNum(ds []rune) (int, error) {
	if len(ds) == 0 {
		return -1, fmt.Errorf("invalid number []runes: %v, string: %v", ds, string(ds))
	}

	m := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
	}

	num := m[ds[len(ds)-1]]

	for i, d := range ds[:len(ds)-1] {
		l := math.Pow(10, float64(len(ds)-i-1))
		v := int(l) * m[d]
		num += v
	}

	return num, nil
}
