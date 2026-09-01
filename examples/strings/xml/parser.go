package xml

import (
	"iter"
	"strconv"
	"unicode"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// entities are the five predefined entities of XML.
var entities = map[string]rune{
	"lt":   '<',
	"gt":   '>',
	"amp":  '&',
	"apos": '\'',
	"quot": '"',
}

// isNameStart reports whether the rune can open an XML name: a
// letter, an underscore, or a colon.
func isNameStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == ':'
}

// isNameRune reports whether the rune can appear inside an XML name.
func isNameRune(r rune) bool {
	return isNameStart(r) || unicode.IsDigit(r) || r == '-' || r == '.'
}

// Events returns the SAX event stream of a document: a combinator
// whose result is a lazy sequence with one event per step. One step
// dispatches on the markup prefix: the longest trie match commits
// the step to its branch, and the branch owns the input from there,
// so its errors reach the stream as they are. The stream ends
// silently at the end of the document and with the error of the
// first malformed event; see [parsec.Seq] for the full contract.
func Events() parsec.Combinator[rune, strings.Position, iter.Seq2[Event, parsec.Error[strings.Position]]] {
	branches := parsec.NewLongestPrefixTree[string, strings.Position, rune, Event](
		map[string]parsec.Combinator[rune, strings.Position, Event]{
			"<!--":      commentBody(),
			"<![CDATA[": cdataBody(),
			"<?":        procInstBody(),
			"</":        endTagBody(),
			"<":         startTagBody(),
		},
		func(s string) []rune {
			return []rune(s)
		},
	)

	text := charData()

	step := func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		pos := buf.Position()

		branch, err := branches.Lookup(buf)
		if err != nil {
			return Event{}, err
		}

		var ev Event

		if branch != nil {
			ev, err = branch(buf)
		} else {
			ev, err = text(buf)
		}

		ev.Pos = pos

		return ev, err
	}

	return strings.Seq(step)
}

// Parse drains the event stream of src and checks the document
// shape: end tags must match open ones, the root element is single,
// and nothing but comments and processing instructions may follow it.
func Parse(src string) ([]Event, error) {
	buf := strings.Buffer([]rune(src))

	seq, err := Events()(buf)
	if err != nil {
		return nil, err
	}

	var events []Event
	var stack []string
	closed := false

	for ev, err := range seq {
		if err != nil {
			return nil, err
		}

		if closed && !allowedAfterRoot(ev) {
			return nil, parsec.NewParseError(ev.Pos, "content after root element")
		}

		switch ev.Kind {
		case KindStartTag:
			stack = append(stack, ev.Name)
		case KindEndTag:
			if len(stack) == 0 || stack[len(stack)-1] != ev.Name {
				return nil, parsec.NewParseError(ev.Pos, "unexpected end tag </"+ev.Name+">")
			}

			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				closed = true
			}
		case KindEmptyTag:
			if len(stack) == 0 {
				closed = true
			}
		}

		events = append(events, ev)
	}

	if len(stack) > 0 {
		return nil, parsec.NewParseError(buf.Position(), "unclosed element <"+stack[len(stack)-1]+">")
	}

	return events, nil
}

// allowedAfterRoot reports whether the event may follow the root
// element: comments and processing instructions always, whitespace
// character data too, everything else is content out of place.
func allowedAfterRoot(ev Event) bool {
	switch ev.Kind {
	case KindComment, KindProcInst:
		return true
	case KindCharData:
		return isSpaceText(ev.Text)
	}

	return false
}

// isSpaceText reports whether the text is whitespace only.
func isSpaceText(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// name parses an XML name: a start rune, then as many name runes as
// there are.
func name() parsec.Combinator[rune, strings.Position, string] {
	head := strings.Try(strings.Satisfy("expected name", true, isNameStart))
	tail := strings.Many(0, strings.Try(strings.Satisfy("expected name", true, isNameRune)))

	return parsec.Cast(
		parsec.And(head, tail, func(head rune, tail []rune) []rune {
			return append([]rune{head}, tail...)
		}),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
}

// commentBody parses the rest of a comment after the dispatch has
// consumed its opener: everything up to -->, which is not included.
func commentBody() parsec.Combinator[rune, strings.Position, Event] {
	body := bodyUntil("-->", "unterminated comment")

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		text, err := body(buf)
		if err != nil {
			return Event{}, err
		}

		return Event{Kind: KindComment, Text: text}, nil
	}
}

// cdataBody parses the rest of a CDATA section after the dispatch
// has consumed its opener: everything up to ]]>, as is — markup
// inside is not interpreted.
func cdataBody() parsec.Combinator[rune, strings.Position, Event] {
	body := bodyUntil("]]>", "unterminated CDATA section")

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		text, err := body(buf)
		if err != nil {
			return Event{}, err
		}

		return Event{Kind: KindCDATA, Text: text}, nil
	}
}

// procInstBody parses the rest of a processing instruction after the
// dispatch has consumed <?: a target name, then everything up to ?>.
// The xml declaration is one, with the target xml.
func procInstBody() parsec.Combinator[rune, strings.Position, Event] {
	target := name()
	body := bodyUntil("?>", "unterminated processing instruction")

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		n, err := target(buf)
		if err != nil {
			return Event{}, parsec.NewParseError(buf.Position(), "expected processing instruction target", err)
		}

		text, err := body(buf)
		if err != nil {
			return Event{}, err
		}

		return Event{Kind: KindProcInst, Name: n, Text: text}, nil
	}
}

// endTagBody parses the rest of a closing tag after the dispatch has
// consumed </: a name, then >.
func endTagBody() parsec.Combinator[rune, strings.Position, Event] {
	target := name()
	gt := strings.Eq("expected '>'", '>')

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		n, err := target(buf)
		if err != nil {
			return Event{}, parsec.NewParseError(buf.Position(), "expected element name", err)
		}

		if _, err := gt(buf); err != nil {
			return Event{}, parsec.NewParseError(buf.Position(), "expected '>' after </"+n, err)
		}

		return Event{Kind: KindEndTag, Name: n}, nil
	}
}

// startTagBody parses the rest of an opening tag after the dispatch
// has consumed <: a name, attributes, then > — or the empty element
// syntax />, which reports the element as a single event. The
// attribute loop is manual: the whitespace before an attribute is
// the probe that stops it, everything after the whitespace is
// committed to the attribute, and a committed body cannot live
// inside Many.
func startTagBody() parsec.Combinator[rune, strings.Position, Event] {
	target := name()
	space := strings.Try(strings.Space("expected attribute"))
	attribute := attr()
	slash := strings.Try(strings.String("expected '>'", "/"))
	gt := strings.Eq("expected '>'", '>')

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		n, err := target(buf)
		if err != nil {
			return Event{}, parsec.NewParseError(buf.Position(), "expected element name", err)
		}

		var as []Attr

		for {
			if _, err := space(buf); err != nil {
				break
			}

			a, err := attribute(buf)
			if err != nil {
				return Event{}, err
			}

			as = append(as, a)
		}

		if _, err := slash(buf); err == nil {
			if _, err := gt(buf); err != nil {
				return Event{}, err
			}

			return Event{Kind: KindEmptyTag, Name: n, Attrs: as}, nil
		}

		if _, err := gt(buf); err != nil {
			return Event{}, err
		}

		return Event{Kind: KindStartTag, Name: n, Attrs: as}, nil
	}
}

// attr parses one attribute after its leading whitespace, which the
// caller has already consumed as the loop probe: a name, =, and a
// quoted value with entities decoded. Everything here is committed:
// a broken attribute is an error, not a missing one.
func attr() parsec.Combinator[rune, strings.Position, Attr] {
	target := name()
	eq := strings.Eq("expected '='", '=')
	value := attrValue()

	return func(buf parsec.Buffer[rune, strings.Position]) (Attr, parsec.Error[strings.Position]) {
		n, err := target(buf)
		if err != nil {
			return Attr{}, err
		}

		if _, err := eq(buf); err != nil {
			return Attr{}, err
		}

		v, err := value(buf)
		if err != nil {
			return Attr{}, err
		}

		return Attr{Name: n, Value: v}, nil
	}
}

// attrValue parses a quoted attribute value with entities decoded.
// The opening quote is a probe; everything after it is committed: a
// '<' inside the value or a missing closing quote is an error, and
// the run stops only at the opening quote again.
func attrValue() parsec.Combinator[rune, strings.Position, string] {
	quote := parsec.Choice("expected quoted value",
		strings.Try(strings.Eq("expected quote", '"')),
		strings.Try(strings.Eq("expected quote", '\'')),
	)

	entity := entityBody()

	return func(buf parsec.Buffer[rune, strings.Position]) (string, parsec.Error[strings.Position]) {
		pos := buf.Position()

		q, err := quote(buf)
		if err != nil {
			return "", err
		}

		body := strings.Try(strings.Satisfy("expected attribute value", true, func(r rune) bool {
			return r != q && r != '<'
		}))

		var out []rune

		for !buf.IsEOF() {
			r, err := body(buf)
			if err != nil {
				break
			}

			if r == '&' {
				decoded, derr := entity(buf)
				if derr != nil {
					return "", derr
				}

				out = append(out, decoded)

				continue
			}

			out = append(out, r)
		}

		if _, err := strings.Eq("expected closing quote", q)(buf); err != nil {
			return "", parsec.NewParseError(pos, "unterminated attribute value", err)
		}

		return string(out), nil
	}
}

// charData parses one run of character data: plain runes up to the
// next '<', with entities and character references decoded. The
// first rune is a probe — the run is committed only once it has
// started, and an empty run fails without consuming input.
func charData() parsec.Combinator[rune, strings.Position, Event] {
	text := strings.Try(strings.Satisfy("expected character data", true, func(r rune) bool {
		return r != '<'
	}))

	entity := entityBody()

	return func(buf parsec.Buffer[rune, strings.Position]) (Event, parsec.Error[strings.Position]) {
		var out []rune

		for !buf.IsEOF() {
			r, err := text(buf)
			if err != nil {
				break
			}

			if r == '&' {
				decoded, derr := entity(buf)
				if derr != nil {
					return Event{}, derr
				}

				out = append(out, decoded)

				continue
			}

			out = append(out, r)
		}

		if len(out) == 0 {
			return Event{}, parsec.NewParseError(buf.Position(), "expected character data")
		}

		return Event{Kind: KindCharData, Text: string(out)}, nil
	}
}

// entityBody parses the body of an entity reference — everything
// between the consumed '&' and the ';': one of the five predefined
// names, a decimal character reference &#65;, or a hexadecimal one
// &#x41;. It is committed: a missing ';' or an unknown name is an
// error, not a literal.
func entityBody() parsec.Combinator[rune, strings.Position, rune] {
	semi := strings.Eq("expected ';'", ';')
	ref := strings.Try(strings.Satisfy("expected entity name", true, func(r rune) bool {
		return isNameRune(r) || r == '#'
	}))

	return func(buf parsec.Buffer[rune, strings.Position]) (rune, parsec.Error[strings.Position]) {
		pos := buf.Position()

		var body []rune

		for !buf.IsEOF() {
			r, err := ref(buf)
			if err != nil {
				break
			}

			body = append(body, r)
		}

		if len(body) == 0 {
			return 0, parsec.NewParseError(pos, "expected entity name after '&'")
		}

		if _, err := semi(buf); err != nil {
			return 0, parsec.NewParseError(pos, "unterminated entity &"+string(body), err)
		}

		return decodeEntity(pos, body)
	}
}

// decodeEntity turns the body of an entity reference into the rune
// it names.
func decodeEntity(pos strings.Position, body []rune) (rune, parsec.Error[strings.Position]) {
	name := string(body)

	if body[0] != '#' {
		if r, ok := entities[name]; ok {
			return r, nil
		}

		return 0, parsec.NewParseError(pos, "unknown entity &"+name+";")
	}

	base := 10
	digits := body[1:]

	if len(digits) > 1 && (digits[0] == 'x' || digits[0] == 'X') {
		base = 16
		digits = digits[1:]
	}

	value, err := strconv.ParseInt(string(digits), base, 32)
	if err != nil || value < 0 {
		return 0, parsec.NewParseError(pos, "malformed character reference &"+name+";")
	}

	return rune(value), nil
}

// bodyUntil parses everything up to the end marker and returns it
// without the marker. It is committed: an exhausted buffer before
// the marker is an error, and a partial match of the marker is just
// text.
func bodyUntil(end string, errMessage string) parsec.Combinator[rune, strings.Position, string] {
	anyRune := strings.Any()

	return func(buf parsec.Buffer[rune, strings.Position]) (string, parsec.Error[strings.Position]) {
		var out []rune
		matched := 0

		for !buf.IsEOF() {
			r, err := anyRune(buf)
			if err != nil {
				return "", err
			}

			if r == rune(end[matched]) {
				matched++

				if matched == len(end) {
					return string(out), nil
				}

				continue
			}

			for i := 0; i < matched; i++ {
				out = append(out, rune(end[i]))
			}

			matched = 0

			if r == rune(end[0]) {
				matched = 1

				continue
			}

			out = append(out, r)
		}

		return "", parsec.NewParseError(buf.Position(), errMessage)
	}
}
