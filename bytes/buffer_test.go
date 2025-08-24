package bytes

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/common"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	t.Parallel()

	type (
		read struct {
			greedy bool
			output byte
			err    error
		}

		seek struct {
			pos int
			err error
		}

		call struct {
			read *read
			seek *seek

			afterPosition int
			afterIsEOF    bool
		}

		test struct {
			input []byte

			beforePosition int
			beforeIsEOF    bool

			calls []call
		}
	)

	tests := []test{
		{
			input:          []byte(""),
			beforePosition: 0,
			beforeIsEOF:    true,
			calls: []call{
				{
					read: &read{
						greedy: false,
						output: 0,
						err:    common.ErrEndOfFile,
					},
					afterPosition: 0,
					afterIsEOF:    true,
				},
				{
					read: &read{
						greedy: true,
						output: 0,
						err:    common.ErrEndOfFile,
					},
					afterPosition: 0,
					afterIsEOF:    true,
				},
				{
					seek: &seek{
						pos: 0,
					},
					afterPosition: 0,
					afterIsEOF:    true,
				},
				{
					seek: &seek{
						pos: 1,
						err: common.ErrOutOfBounds,
					},
					afterPosition: 0,
					afterIsEOF:    true,
				},
				{
					seek: &seek{
						pos: -1,
						err: common.ErrOutOfBounds,
					},
					afterPosition: 0,
					afterIsEOF:    true,
				},
			},
		},
		{
			input:          []byte("foo"),
			beforePosition: 0,
			beforeIsEOF:    false,
			calls: []call{
				{
					read: &read{
						greedy: false,
						output: 'f',
						err:    nil,
					},
					afterPosition: 0,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'f',
						err:    nil,
					},
					afterPosition: 1,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: false,
						output: 'o',
						err:    nil,
					},
					afterPosition: 1,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
						err:    nil,
					},
					afterPosition: 2,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: false,
						output: 'o',
						err:    nil,
					},
					afterPosition: 2,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
						err:    nil,
					},
					afterPosition: 3,
					afterIsEOF:    true,
				},
			},
		},
		{
			input:          []byte("foo"),
			beforePosition: 0,
			beforeIsEOF:    false,
			calls: []call{
				{
					read: &read{
						greedy: false,
						output: 'f',
						err:    nil,
					},
					afterPosition: 0,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'f',
						err:    nil,
					},
					afterPosition: 1,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: false,
						output: 'o',
						err:    nil,
					},
					afterPosition: 1,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
						err:    nil,
					},
					afterPosition: 2,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: false,
						output: 'o',
						err:    nil,
					},
					afterPosition: 2,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
						err:    nil,
					},
					afterPosition: 3,
					afterIsEOF:    true,
				},
				{
					seek: &seek{
						pos: 100,
						err: common.ErrOutOfBounds,
					},
					afterPosition: 3,
					afterIsEOF:    true,
				},
				{
					seek: &seek{
						pos: 2,
					},
					afterPosition: 2,
					afterIsEOF:    false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
						err:    nil,
					},
					afterPosition: 3,
					afterIsEOF:    true,
				},
				{
					read: &read{
						greedy: true,
						output: 0,
						err:    common.ErrEndOfFile,
					},
					afterPosition: 3,
					afterIsEOF:    true,
				},
			},
		},
	}

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("case %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := Buffer(test.input)

			assert.Equal(t, b.Position(), example.beforePosition)
			assert.Equal(t, b.IsEOF(), example.beforeIsEOF)

			for i, call := range test.calls {
				t.Logf("call %d", i)

				if call.read != nil {
					result, err := b.Read(call.read.greedy)

					if call.read.err == nil {
						assert.NoError(t, err)
					} else {
						assert.Error(t, err)
						assert.EqualError(t, err, call.read.err.Error())
					}

					assert.Equal(t, result, call.read.output)
				} else if call.seek != nil {
					err := b.Seek(call.seek.pos)

					if call.seek.err == nil {
						assert.NoError(t, err)
					} else {
						assert.Error(t, err)
						assert.EqualError(t, err, call.seek.err.Error())
					}
				} else {
					t.Fatal("invalid test")
				}

				assert.Equal(t, b.Position(), call.afterPosition)
				assert.Equal(t, b.IsEOF(), call.afterIsEOF)
			}
		})
	}
}
