package xml

import (
	"fmt"

	"github.com/okneniz/parsec/strings"
)

// Kind is the type of one SAX event.
type Kind uint8

// The event kinds, in the order they can appear in a document.
const (
	KindProcInst Kind = iota
	KindComment
	KindCDATA
	KindStartTag
	KindEndTag
	KindEmptyTag
	KindCharData
)

// String returns the name of the kind.
func (k Kind) String() string {
	switch k {
	case KindProcInst:
		return "procinst"
	case KindComment:
		return "comment"
	case KindCDATA:
		return "cdata"
	case KindStartTag:
		return "start"
	case KindEndTag:
		return "end"
	case KindEmptyTag:
		return "empty"
	case KindCharData:
		return "text"
	}

	return "unknown"
}

// Attr is one attribute of a start tag: a name and a decoded value.
type Attr struct {
	Name  string
	Value string
}

// Event is one SAX event: where it starts, what it is, and its parts.
// A start or empty tag carries the element name and its attributes,
// an end tag the name, a processing instruction the target and the
// rest of the text, a comment, CDATA section, or character data the
// text itself.
type Event struct {
	Kind  Kind
	Name  string
	Text  string
	Attrs []Attr
	Pos   strings.Position
}

// String renders the event on one line: the kind, then its parts.
func (e Event) String() string {
	switch e.Kind {
	case KindStartTag, KindEmptyTag:
		line := fmt.Sprintf("%s %s", e.Kind, e.Name)

		for _, attr := range e.Attrs {
			line += fmt.Sprintf(" %s=%q", attr.Name, attr.Value)
		}

		return line
	case KindEndTag:
		return fmt.Sprintf("%s %s", e.Kind, e.Name)
	case KindProcInst:
		return fmt.Sprintf("%s %s %q", e.Kind, e.Name, e.Text)
	}

	return fmt.Sprintf("%s %q", e.Kind, e.Text)
}
