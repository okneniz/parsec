// Package xml is a SAX parser over a rune buffer: the document is a
// lazy stream of events, one combinator step per event, and the
// consumer sees the document as it is being read — a break in the
// range leaves the buffer right after the last event taken.
//
// It is a study example of the Seq combinator, not a full XML
// implementation: no namespaces, DTD, or character encodings. The
// supported surface is processing instructions, comments, CDATA
// sections, elements with quoted attributes, character data with the
// five predefined entities and character references, and the empty
// element syntax <br/>, reported as a single event.
package xml
