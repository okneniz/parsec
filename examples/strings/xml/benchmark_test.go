package xml

import (
	"testing"

	"github.com/okneniz/parsec/strings"
)

// sink keeps the consumed events from being optimized away.
var sink int

// BenchmarkSAX consumes an endless document: the buffer holds one
// cell of five events, and when the cell is exhausted the benchmark
// rewinds the buffer and the same stream continues from the start.
// One op is one event, so with -benchmem the B/op column is the
// memory cost of a single event — and it stays the same however many
// events the framework asks for: the parser holds nothing but the
// buffer and the event in hand.
func BenchmarkSAX(b *testing.B) {
	buf := strings.Buffer([]rune(`<a x="1">text &amp; more</a><!-- note --><b/>`))
	start := buf.Position()

	seq, _ := Events()(buf)

	b.ReportAllocs()

	for b.Loop() {
		if buf.IsEOF() {
			if err := buf.Seek(start); err != nil {
				b.Fatal(err)
			}
		}

		for ev, err := range seq {
			if err != nil {
				b.Fatal(err)
			}

			sink += len(ev.Name) + len(ev.Text)

			break
		}
	}
}
