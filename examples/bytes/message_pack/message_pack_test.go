package message_pack

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	b "git.sr.ht/~okneniz/parsec/bytes"
	// c "git.sr.ht/~okneniz/parsec/common"

	mpack "github.com/vmihailenco/msgpack/v5"
)

func TestMain_MessagePack(t *testing.T) {
	t.Run("primitives", func(t *testing.T){
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
			"test                           asd                  asd", // long string
			nil,
			-1,
			0,
			1,
		}

		for _, example := range examples {
			expected := example
			name := fmt.Sprintf("%#T - %v", example, example)

			t.Run(name, func(t *testing.T) {
				data, err := mpack.Marshal(expected)
				if err != nil {
					t.Fatal(err)
				}

				actual, err := b.Parse(data, MessagePack())
				if err != nil {
					t.Errorf("input: %v", data)
					t.Errorf("hex: %s", toHex(data))
					t.Errorf("bin: %s", toBinary(data))
					t.Fatal(err)
				}

				if reflect.DeepEqual(expected, actual) {
					t.Errorf("expected: %v", expected)
					t.Errorf("actual: %v", actual)
					t.Fatal()
				}
			})
		}

		// for i := 0; i > - 127; i-- {
		// 	t.Log("input", i)
		// 	data, err := mpack.Marshal(i)
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	t.Log("output", data)
		// }
	})
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
