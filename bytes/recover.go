package bytes

import (
	"github.com/okneniz/parsec"
)

// Try applies c and, when it fails, rewinds the buffer to the previous
// position. This is the backtracking primitive of the library:
// greedy combinators consume input even when they fail, so every
// combinator which can fail inside an alternative branch (see Choice,
// Or) or in a loop must be wrapped in Try to make backtracking possible.
// The state is not rewound: a combinator that changes it must restore
// the change itself before failing.
func Try[O any, S parsec.Recoverable[S]](
	c parsec.Combinator[byte, int, O, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Try[byte, int, O, S](c)
}
