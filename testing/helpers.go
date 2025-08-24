package testing

import (
	"testing"
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
