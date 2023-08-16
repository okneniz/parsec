package message_pack

import (
	"fmt"
	"math"
	"reflect"
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
			nil,
		}

		for _, example := range examples {
			expected := example
			name := fmt.Sprintf("%#T", example)

			t.Run(name, func(t *testing.T) {
				data, err := mpack.Marshal(expected)
				if err != nil {
					t.Fatal(err)
				}

				actual, err := b.Parse(data, MessagePack())
				if err != nil {
					t.Errorf("input: %v", data)
					t.Fatal(err)
				}

				if reflect.DeepEqual(expected, actual) {
					t.Errorf("expected: %v", expected)
					t.Errorf("actual: %v", actual)
					t.Fatal()
				}
			})
		}
	})
}
