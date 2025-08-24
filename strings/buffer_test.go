package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okneniz/parsec/common"
)

func TestBuffer(t *testing.T) {
	t.Parallel()

	type (
		read struct {
			greedy bool
			output rune
			err    error
		}

		seek struct {
			pos Position
			err error
		}

		call struct {
			read *read
			seek *seek

			afterPosition Position
			afterIsEOF    bool
		}

		test struct {
			input string

			beforePosition Position
			beforeIsEOF    bool

			calls []call
		}
	)

	tests := []test{
		{
			input: "",
			beforePosition: Position{
				line:   0,
				column: 0,
				index:  0,
			},
			beforeIsEOF: true,
			calls: []call{
				{
					read: &read{
						greedy: false,
						output: 0,
						err:    common.ErrEndOfFile,
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
				{
					read: &read{
						greedy: true,
						output: 0,
						err:    common.ErrEndOfFile,
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
				{
					seek: &seek{
						pos: Position{
							line:   0,
							column: 0,
							index:  0,
						},
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
				{
					seek: &seek{
						pos: Position{
							line:   0,
							column: 1,
							index:  1,
						},
						err: common.ErrOutOfBounds,
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
				{
					seek: &seek{
						pos: Position{
							line:   0,
							column: 0,
							index:  -1,
						},
						err: common.ErrOutOfBounds,
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
				{
					seek: &seek{
						pos: Position{
							line:   0,
							column: 0,
							index:  0,
						},
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: true,
				},
			},
		},
		{
			input: "foo\nbar\nbaz",
			beforePosition: Position{
				line:   0,
				column: 0,
				index:  0,
			},
			beforeIsEOF: false,
			calls: []call{
				{
					read: &read{
						greedy: false,
						output: 'f',
					},
					afterPosition: Position{
						line:   0,
						column: 0,
						index:  0,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'f',
					},
					afterPosition: Position{
						line:   0,
						column: 1,
						index:  1,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
					},
					afterPosition: Position{
						line:   0,
						column: 2,
						index:  2,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'o',
					},
					afterPosition: Position{
						line:   0,
						column: 3,
						index:  3,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: '\n',
					},
					afterPosition: Position{
						line:   1,
						column: 0,
						index:  4,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'b',
					},
					afterPosition: Position{
						line:   1,
						column: 1,
						index:  5,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'a',
					},
					afterPosition: Position{
						line:   1,
						column: 2,
						index:  6,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'r',
					},
					afterPosition: Position{
						line:   1,
						column: 3,
						index:  7,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: '\n',
					},
					afterPosition: Position{
						line:   2,
						column: 0,
						index:  8,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'b',
					},
					afterPosition: Position{
						line:   2,
						column: 1,
						index:  9,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'a',
					},
					afterPosition: Position{
						line:   2,
						column: 2,
						index:  10,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'z',
					},
					afterPosition: Position{
						line:   2,
						column: 3,
						index:  11,
					},
					afterIsEOF: true,
				},
				{
					seek: &seek{
						pos: Position{
							line:   2,
							column: 1,
							index:  9,
						},
					},
					afterPosition: Position{
						line:   2,
						column: 1,
						index:  9,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'a',
					},
					afterPosition: Position{
						line:   2,
						column: 2,
						index:  10,
					},
					afterIsEOF: false,
				},
				{
					read: &read{
						greedy: true,
						output: 'z',
					},
					afterPosition: Position{
						line:   2,
						column: 3,
						index:  11,
					},
					afterIsEOF: true,
				},
			},
		},
	}

	for i, example := range tests {
		test := example
		name := fmt.Sprintf("case %d", i)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := Buffer([]rune(test.input))

			assert.Equal(t, b.Position(), example.beforePosition)
			assert.Equal(t, b.IsEOF(), example.beforeIsEOF)

			for i, call := range test.calls {
				t.Logf("call %d", i)

				if call.read != nil {
					result, err := b.Read(call.read.greedy)

					if call.read.err == nil {
						assert.NoError(t, err)
					} else {
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

				assert.Equal(t, b.Position().String(), call.afterPosition.String())
				assert.Equal(t, b.IsEOF(), call.afterIsEOF)
			}
		})
	}
}
