package bytes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/assert"

	"github.com/okneniz/parsec/common"
)

type (
	testCase[T any] struct {
		input  []byte
		output T
		err    error
	}

	test[T any] struct {
		comb  common.Combinator[byte, int, T]
		cases []testCase[T]
	}
)

func runTests[T comparable](t *testing.T, tests []test[T]) {
	t.Helper()

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("test %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, x := range test.cases {
				testCase := x
				name := fmt.Sprintf("case %d", i)

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					result, err := Parse(testCase.input, test.comb)

					if testCase.err != nil {
						assert.Error(t, err)
						assert.Equal(t, err.Error(), testCase.err.Error())
					} else {
						assert.NoError(t, err)
					}

					assert.Equal(t, testCase.output, result)
				})
			}
		})
	}
}

func runTestsSlice[T comparable](t *testing.T, tests []test[[]T]) {
	t.Helper()

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("test %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for i, x := range test.cases {
				testCase := x
				name := fmt.Sprintf("case %d", i)

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					result, err := Parse(testCase.input, test.comb)

					if testCase.err != nil {
						assert.Error(t, err)
						assert.EqualError(t, err, testCase.err.Error())
					} else {
						assert.NoError(t, err)
					}

					assert.EqualValues(t, testCase.output, result)
				})
			}
		})
	}
}

func TestSatisfy(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Satisfy("expected not 'c'", true, func(x byte) bool { return x != 'c' }),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "expected not 'c'"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("b"),
					output: 'b',
				},
				{
					input:  []byte("c"),
					output: 0,
					err:    common.NewParseError(0, "expected not 'c'"),
				},
			},
		},
		{
			comb: Satisfy("error explanation", true, common.Nothing),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
				{
					input:  []byte("abc"),
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
			},
		},
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Any(),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, common.ErrEndOfFile.Error()),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
			},
		},
	})

	t.Run("must parse first element from none empty input", func(t *testing.T) {
		t.Parallel()

		seed := time.Now().UnixNano()
		t.Logf("seed: %v", seed)
		rnd := rand.New(rand.NewPCG(0, uint64(seed)))
		comb := Any()

		ohsnap.Check(t, 1000, ohsnap.NewBuilder(rnd).Byte(), func(b byte) bool {
			result, err := Parse([]byte{b}, comb)
			if err != nil {
				t.Logf("input: %v", []byte{b})
				t.Logf("output: %v", result)
				t.Error(err)
				return false
			}

			return result == b
		})
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	runTests(t, []test[byte]{
		{
			comb: Try(
				Satisfy(
					"error explanation",
					true,
					func(x byte) bool { return x <= byte('b') },
				),
			),
			cases: []testCase[byte]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "error explanation"),
				},
				{
					input:  []byte("a"),
					output: 'a',
				},
				{
					input:  []byte("b"),
					output: 'b',
				},
				{
					input: []byte("c"),
					err:   common.NewParseError(0, "error explanation"),
				},
			},
		},
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	notBrackets := Satisfy("test", true, func(x byte) bool {
		return !(x == byte(')') || x == byte('('))
	})

	comb := Between(
		Eq("expected '('", '('),
		Some(0, "expected not ( or ) symbols", Try(notBrackets)),
		Eq("expected ')'", ')'),
	)

	runTestsSlice(t, []test[[]byte]{
		{
			comb: comb,
			cases: []testCase[[]byte]{
				{
					input:  []byte{},
					output: nil,
					err:    common.NewParseError(0, "expected '('"),
				},
				{
					input:  []byte("(abc)"),
					output: []byte{'a', 'b', 'c'},
				},
				{
					input: []byte("abc"),
					err:   common.NewParseError(0, "expected '('"),
				},
				{
					input: []byte("(abc"),
					err:   common.NewParseError(4, "expected ')'"),
				},
				{
					input: []byte(" (abc) "),
					err:   common.NewParseError(0, "expected '('"),
				},
				{
					input: []byte("((abc))"),
					err:   common.NewParseError(1, "expected not ( or ) symbols"),
				},
				{
					input: []byte("()"),
					err:   common.NewParseError(1, "expected not ( or ) symbols"),
				},
			},
		},
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	runTests(t, []test[bool]{
		{
			comb: EOF(),
			cases: []testCase[bool]{
				{
					input:  []byte{},
					output: true,
				},
				{
					input:  []byte("1"),
					output: false,
				},
				{
					input:  []byte("123"),
					output: false,
				},
			},
		},
	})
}

func TestCast(t *testing.T) {
	t.Parallel()

	runTests(t, []test[int]{
		{
			comb: Cast(
				Satisfy("test", true, common.Anything[byte]),
				func(x byte) (int, error) { return int(x), nil },
			),
			cases: []testCase[int]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, "test"),
				},
				{
					input:  []byte{0},
					output: 0,
				},
				{
					input:  []byte{1},
					output: 1,
				},
				{
					input:  []byte{math.MaxUint8},
					output: 255,
				},
				{
					input:  []byte{math.MaxInt8},
					output: 127,
				},
			},
		},
		{
			comb: Cast(
				Any(),
				func(x byte) (int, error) { return -1, fmt.Errorf("test error") },
			),
			cases: []testCase[int]{
				{
					input:  []byte{},
					output: 0,
					err:    common.NewParseError(0, common.ErrEndOfFile.Error()),
				},
				{
					input:  []byte{0},
					output: -1,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{1},
					output: -1,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{math.MaxUint8},
					output: -1,
					err:    common.NewParseError(0, "test error"),
				},
				{
					input:  []byte{math.MaxInt8},
					output: -1,
					err:    common.NewParseError(0, "test error"),
				},
			},
		},
	})
}

func TestReadAs(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	rnd := rand.New(rand.NewPCG(0, uint64(seed)))

	t.Run("uint8", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint8(),
				binary.BigEndian,
				ReadAs[uint8](1, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint8(),
				binary.LittleEndian,
				ReadAs[uint8](1, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("uint16", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint16(),
				binary.BigEndian,
				ReadAs[uint16](2, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint16(),
				binary.LittleEndian,
				ReadAs[uint16](2, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("uint32", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint32(),
				binary.BigEndian,
				ReadAs[uint32](4, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint32(),
				binary.LittleEndian,
				ReadAs[uint32](4, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint64(),
				binary.BigEndian,
				ReadAs[uint64](8, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Uint64(),
				binary.LittleEndian,
				ReadAs[uint64](8, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("int8", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int8(),
				binary.BigEndian,
				ReadAs[int8](1, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int8(),
				binary.LittleEndian,
				ReadAs[int8](1, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("int16", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int16(),
				binary.BigEndian,
				ReadAs[int16](2, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int16(),
				binary.LittleEndian,
				ReadAs[int16](2, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("int32", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int32(),
				binary.BigEndian,
				ReadAs[int32](4, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int32(),
				binary.LittleEndian,
				ReadAs[int32](4, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int64(),
				binary.BigEndian,
				ReadAs[int64](8, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Int64(),
				binary.LittleEndian,
				ReadAs[int64](8, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("float32", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Float32(),
				binary.BigEndian,
				ReadAs[float32](4, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Float32(),
				binary.LittleEndian,
				ReadAs[float32](4, "E", binary.LittleEndian),
			)
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()

		t.Run("big endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Float64(),
				binary.BigEndian,
				ReadAs[float64](8, "E", binary.BigEndian),
			)
		})

		t.Run("little endian", func(t *testing.T) {
			t.Parallel()

			checkReadAs(
				t,
				ohsnap.NewBuilder(rnd).Float64(),
				binary.LittleEndian,
				ReadAs[float64](8, "E", binary.LittleEndian),
			)
		})
	})
}

func checkReadAs[T comparable](
	t *testing.T,
	arbT ohsnap.Arbitrary[T],
	order binary.ByteOrder,
	comb common.Combinator[byte, int, T],
) {
	t.Helper()

	const iterations = 100_000

	arb := ohsnap.Map(arbT, func(input T) ohsnap.Pair[T, []byte] {
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, order, &input)

		return ohsnap.Pair[T, []byte]{
			First:  input,
			Second: buf.Bytes(),
		}
	},
	)

	ohsnap.Check(t, iterations, arb, func(data ohsnap.Pair[T, []byte]) bool {
		expected := data.First
		input := data.Second

		actual, err := Parse(input, comb)
		if err != nil {
			t.Logf("input: %v", input)
			t.Logf("output: %v", actual)
			t.Error(err)
			return false
		}

		return assert.EqualValues(t, expected, actual)
	})
}
