package xml_test

import (
	"fmt"

	"github.com/okneniz/parsec/examples/strings/xml"
	"github.com/okneniz/parsec/strings"
)

// Events is lazy: the range reads the document event by event, and
// the break stops it mid-document — the buffer stays right after the
// last event taken, so a second stream picks up where the first one
// stopped.
func ExampleEvents() {
	buf := strings.Buffer([]rune(`<library><book id="1"/><book id="2"/></library>`))

	seq, _ := xml.Events()(buf)

	taken := 0

	for ev, err := range seq {
		if err != nil {
			fmt.Println(err)

			break
		}

		fmt.Println(ev)

		taken++

		if taken == 3 {
			break
		}
	}

	rest, _ := xml.Events()(buf)

	for ev, err := range rest {
		if err != nil {
			fmt.Println(err)

			break
		}

		fmt.Println(ev)
	}

	// Output:
	// start library
	// empty book id="1"
	// empty book id="2"
	// end library
}

// Parse drains the whole stream and checks the shape of the
// document: tags nest, the root is single, nothing follows it.
func ExampleParse() {
	events, err := xml.Parse(`<?xml version="1.0"?><note to="you">&amp;done</note>`)
	fmt.Println(err)

	for _, ev := range events {
		fmt.Println(ev)
	}

	// Output:
	// <nil>
	// procinst xml " version=\"1.0\""
	// start note to="you"
	// text "&done"
	// end note
}
