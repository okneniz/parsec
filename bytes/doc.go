// Package bytes provides parser combinators for binary input:
// a byte Buffer (also constructible from a file), fixed-size
// big/little-endian readers, and thin typed wrappers around
// the generic combinators of [github.com/okneniz/parsec].
//
// Parsing starts from [Parse] (or [ParseFile]), a combinator and
// the parser state, parsec.Stateless when the parser keeps no
// state:
//
//	be := bytes.ReadAs[uint16, parsec.Stateless](2, "expected uint16", binary.BigEndian)
//
//	result, err := bytes.Parse(parsec.Stateless{}, []byte{0x01, 0x02}, be)
//	// result == 258
package bytes
