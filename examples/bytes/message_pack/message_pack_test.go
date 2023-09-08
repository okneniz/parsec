package message_pack

import (
	"fmt"
	"math"
	// "reflect"
	"strings"
	"testing"
	"math/rand"
	"time"

	b "git.sr.ht/~okneniz/parsec/bytes"
	// c "git.sr.ht/~okneniz/parsec/common"

	mpack "github.com/vmihailenco/msgpack/v5"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
			randStringRunes(1),
			randStringRunes(10),
			randStringRunes(255),
			randStringRunes(256),
			randStringRunes(65535),
			randStringRunes(65536),
			[]interface{}{nil, true, false, 1,2,3,"test"},
		}

		for _, example := range examples {
			expected := example
			name := fmt.Sprintf("%#T - %v", example, example)

			t.Run(name, func(t *testing.T) {
				data, err := mpack.Marshal(expected)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("input: %v", data)
				t.Logf("hex: %s", toHex(data))
				t.Logf("bin: %s", toBinary(data))

				actual, err := b.Parse(data, MessagePack())
				if err != nil {
					t.Errorf("output: %v", actual)
					t.Fatal(err)
				}

				if fmt.Sprintf("%v", expected) != actual.String() {
					t.Errorf("expected: %v", expected)
					t.Errorf("actual: %s", actual.String())
					t.Fatal()
				}
			})
		}
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
