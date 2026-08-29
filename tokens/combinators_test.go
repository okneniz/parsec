package tokens

import (
	stdstrings "strings"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// charLexer skips spaces and takes one rune as a char token, following
// the Lexer contract: trivia skipped inside, the exhausted input
// reported with parsec.ErrEndOfFile.
func charLexer(buf parsec.Buffer[rune, strings.Position]) (Token[string, string], error) {
	for !buf.IsEOF() {
		r, err := buf.Read(false)
		if err != nil || r != ' ' {
			break
		}

		if _, err := buf.Read(true); err != nil {
			return Token[string, string]{}, err
		}
	}

	if buf.IsEOF() {
		return Token[string, string]{}, parsec.ErrEndOfFile
	}

	r, err := buf.Read(true)
	if err != nil {
		return Token[string, string]{}, err
	}

	return Token[string, string]{Kind: "char", Lexeme: string(r)}, nil
}

func TestSatisfy(t *testing.T) {
	buf := strings.Buffer([]rune("ab c"))
	exact := Exact(charLexer, "char", "a")
	ofKind := OfKind(charLexer, "char")

	if _, err := exact(buf); err != nil {
		t.Fatalf("Exact a: %v", err)
	}

	if _, err := ofKind(buf); err != nil {
		t.Fatalf("OfKind char: %v", err)
	}

	// greedy failure consumes the token: the c is gone now
	if _, err := Exact(charLexer, "char", "z")(buf); err == nil {
		t.Fatal("Exact char z matched")
	}

	if pos := buf.Position(); pos.Column() != 4 {
		t.Fatalf("position after greedy failure = %v", pos)
	}
}

func TestUnexpected(t *testing.T) {
	buf := strings.Buffer([]rune("ab"))

	unexpected := parsec.Unexpected(charLexer, "testing")
	if _, uerr := unexpected(buf); uerr == nil || !stdstrings.Contains(uerr.Error(), `char "a"`) {
		t.Fatalf("Unexpected = %v", uerr)
	}

	empty := strings.Buffer(nil)

	_, uerr := unexpected(empty)
	if uerr == nil || !stdstrings.Contains(uerr.Error(), "unexpected end of input") {
		t.Fatalf("Unexpected at EOF = %v", uerr)
	}
}

func TestLineComment(t *testing.T) {
	line := LineComment("//")

	// the comment runs to and including the newline
	buf := strings.Buffer([]rune("// hello\nx"))

	ok, err := line(buf)
	if err != nil || !ok {
		t.Fatalf("LineComment failed: %v", err)
	}

	if pos := buf.Position(); pos.Column() != 0 || pos.Line() != 1 {
		t.Fatalf("position after the comment = %v, want line 1 column 0", pos)
	}

	// a comment at the end of input closes at the end
	ok, err = line(strings.Buffer([]rune("// tail")))
	if err != nil || !ok {
		t.Fatalf("LineComment at EOF failed: %v", err)
	}

	// a missing marker fails with the input intact
	buf = strings.Buffer([]rune("x"))

	if _, err := line(buf); err == nil {
		t.Fatal("LineComment of x succeeded, want error")
	}

	if pos := buf.Position(); pos.Column() != 0 {
		t.Fatalf("position after the miss = %v, want column 0", pos)
	}
}

func TestCommentBody(t *testing.T) {
	body := CommentBody("(*", "*)", true)
	open := strings.Try(strings.String("comment start", "(*"))

	// the body runs to the matching close, nested comments counted
	buf := strings.Buffer([]rune("(* a (* b *) c *)+"))

	if _, err := open(buf); err != nil {
		t.Fatalf("open failed: %v", err)
	}

	ok, err := body(buf)
	if err != nil || !ok {
		t.Fatalf("CommentBody failed: %v", err)
	}

	r, rerr := buf.Read(false)
	if rerr != nil || r != '+' {
		t.Fatalf("rune after the comment = %v %v, want +", r, rerr)
	}

	// an unterminated comment fails
	if _, err := body(strings.Buffer([]rune("(* a"))); err == nil || !stdstrings.Contains(err.Error(), "unterminated comment") {
		t.Fatalf("CommentBody of an unterminated comment = %v", err)
	}

	// with nesting off, an inner open is plain text
	flat := CommentBody("(*", "*)", false)

	buf = strings.Buffer([]rune("(* a *)"))

	if _, err := open(buf); err != nil {
		t.Fatalf("open failed: %v", err)
	}

	if ok, err := flat(buf); err != nil || !ok {
		t.Fatalf("flat CommentBody failed: %v", err)
	}
}
