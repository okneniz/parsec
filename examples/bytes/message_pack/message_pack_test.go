package message_pack

import (
	"fmt"
	"math"

	"math/rand"
	"strings"
	"testing"
	"time"

	b "github.com/okneniz/parsec/bytes"
	"golang.org/x/exp/constraints"

	mpack "github.com/vmihailenco/msgpack/v5"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain_MessagePack(t *testing.T) {
	t.Parallel()

	t.Run("fixed", func(t *testing.T) {
		examples := []interface{}{
			math.MaxFloat32,
			math.MaxFloat64,
			math.MaxInt,
			math.MaxInt8,
			math.MaxInt16,
			math.MaxInt32,
			math.MaxInt64,
			uint(math.MaxUint),
			math.MaxUint8,
			math.MaxUint16,
			math.MaxUint32,
			uint64(math.MaxUint64),
			"test",
			"test                           asd                  asd",
			nil,
			-1,
			-2,
			-30,
			-31,
			-32,
			-33,
			0,
			1,
			[]interface{}{nil, true, false, 1, 2, 3, "test"},
		}

		for _, expected := range examples {
			name := fmt.Sprintf("%#T - %v", expected, expected)

			t.Run(name, func(t *testing.T) {
				check[any](t, expected)
			})
		}
	})

	t.Run("random", func(t *testing.T) {
		examplesCount := 10000

		t.Run("string", func(t *testing.T) {
			t.Run("fixstr", func(t *testing.T) {
				for i := 0; i < examplesCount; i++ {
					size := rnd[int](0, 31)
					expected := randStringRunes(size)
					check[String](t, expected)
				}
			})

			t.Run("str 8", func(t *testing.T) {
				for i := 0; i < examplesCount; i++ {
					size := rnd[int](32, math.Pow(2, 8)-1)
					expected := randStringRunes(size)
					check[String](t, expected)
				}
			})

			t.Run("str 16", func(t *testing.T) {
				for i := 0; i < examplesCount; i++ {
					size := rnd[int](math.Pow(2, 8), math.Pow(2, 16)-1)
					expected := randStringRunes(size)
					check[String](t, expected)
				}
			})

			// t.Run("str 32", func(t *testing.T) {
			// 	for i := 0; i < 1; i++ {
			// 		size := rnd(int(math.Pow(2, 16)), int(math.Pow(2, 32)) - 1)
			// 		expected := randStringRunes(size)
			// 		check[String](t, expected)
			// 	}
			// })
		})

		t.Run("numbers", func(t *testing.T) {
			t.Run("signed", func(t *testing.T) {
				t.Run("negative fixint", func(t *testing.T) {
					for i := -1; i > -31; i-- {
						check[Signed8](t, i)
					}
				})

				t.Run("int 16", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[int16](math.MinInt16, math.MaxInt16)
						check[Signed16](t, number)
						check[Signed16](t, -number)
					}
				})

				t.Run("int 32", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[int32](math.MinInt32, math.MaxInt32)
						check[Signed32](t, number)
						check[Signed32](t, -number)
					}
				})

				t.Run("int 64", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[int64](math.MinInt64, math.MaxInt64)
						check[Signed64](t, number)
						check[Signed64](t, -number)
					}
				})

				t.Run("float 32", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rand.Float32()
						check[Float32](t, number)
						check[Float32](t, -number)
					}
				})

				t.Run("float 64", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rand.Float64()
						check[Float64](t, number)
						check[Float64](t, -number)
					}
				})
			})

			t.Run("unsigned", func(t *testing.T) {
				t.Run("positive fixint", func(t *testing.T) {
					for i := 0; i < int(math.Pow(2, 7)-1); i++ {
						check[Unsigned8](t, i)
					}
				})

				t.Run("uint 16", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[uint16](0, math.MaxInt16)
						check[Unsigned16](t, number)
					}
				})

				t.Run("uint 32", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[uint32](0, math.MaxUint32)
						check[Unsigned32](t, number)
					}
				})

				t.Run("uint 64", func(t *testing.T) {
					for i := 0; i < examplesCount; i++ {
						number := rnd[uint64](0, math.MaxUint64)
						check[Unsigned64](t, number)
					}
				})
			})
		})

		t.Run("binary", func(t *testing.T) {
			t.Run("bin 8", func(t *testing.T) {
				max := math.Pow(2, 8) - 1

				for i := 0; i < examplesCount; i++ {
					size := rnd[int](0, max)
					data := randBytes(size)
					check[Binary](t, data)
				}
			})

			t.Run("bin 16", func(t *testing.T) {
				min := math.Pow(2, 8)
				max := math.Pow(2, 16) - 1

				for i := 0; i < examplesCount; i++ {
					size := rnd[int](min, max)
					data := randBytes(size)
					check[Binary](t, data)
				}
			})

			// t.Run("bin 32", func(t *testing.T) {
			// 	min := math.Pow(2, 16)
			// 	max := math.Pow(2, 32) - 1

			// 	for i := 0; i < 1; i++ {
			// 		size := rnd[int](min, max)
			// 		data := randBytes(size)
			// 		check[Binary](t, data)
			// 	}
			// })
		})
	})
}

func rnd[T constraints.Integer](min, max float64) T {
	return T(rand.Intn(int(max - min + 1 + min)))
}

func check[T any](t *testing.T, expected interface{}) {
	t.Helper()

	data, err := mpack.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := b.Parse(data, MessagePack())
	if err != nil {
		t.Logf("input: %v", data)
		t.Logf("hex: %s", toHex(data))
		t.Logf("bin: %s", toBinary(data))
		t.Errorf("output: %v", actual)
		t.Fatal(err)
	}

	if fmt.Sprintf("%v", expected) != actual.String() {
		t.Logf("input: %v", data)
		t.Logf("hex: %s", toHex(data))
		t.Logf("bin: %s", toBinary(data))
		t.Errorf("output: %v", actual)

		t.Errorf("expected: %v", expected)
		t.Errorf("actual: %s", actual.String())
		t.Fatal()
	}

	if _, ok := actual.(T); !ok {
		t.Logf("unexpected type %T %v", actual, actual)
	}
}

func toHex(data []byte) string {
	result := make([]string, len(data))

	for i, x := range data {
		result[i] = fmt.Sprintf("0x%x", x)
	}

	return "[" + strings.Join(result, " ") + "]"
}

func toBinary(data []byte) string {
	result := make([]string, len(data))

	for i, x := range data {
		result[i] = fmt.Sprintf("%08b", x)
	}

	return "[" + strings.Join(result, " ") + "]"
}

func randBytes(count int) []byte {
	data := make([]byte, count)
	for i := 0; i < count; i++ {
		data[i] = byte(rnd[uint8](0, math.MaxUint8))
	}
	return data
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
