package tml

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		kinds []Kind
		want  []string
	}{
		{
			name:  "keywords are not identifiers",
			src:   "let fn fun if then else end in val true false iff",
			kinds: []Kind{KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindKeyword, KindIdent},
			want:  []string{"let", "fn", "fun", "if", "then", "else", "end", "in", "val", "true", "false", "iff"},
		},
		{
			name:  "maximal munch of symbolic runs",
			src:   "=> = <> <= + - * / ~",
			kinds: []Kind{KindOperator, KindOperator, KindOperator, KindOperator, KindOperator, KindOperator, KindOperator, KindOperator, KindOperator},
			want:  []string{"=>", "=", "<>", "<=", "+", "-", "*", "/", "~"},
		},
		{
			name:  "integers and punctuation",
			src:   "42 0 (x, _)",
			kinds: []Kind{KindInt, KindInt, KindSymbol, KindIdent, KindSymbol, KindSymbol, KindSymbol},
			want:  []string{"42", "0", "(", "x", ",", "_", ")"},
		},
		{
			name:  "nested comments are skipped",
			src:   "1 (* a (* b (* c *) d *) e *) 2",
			kinds: []Kind{KindInt, KindInt},
			want:  []string{"1", "2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := Lex(test.src)
			if err != nil {
				t.Fatalf("Lex(%q) failed: %v", test.src, err)
			}

			if len(tokens) != len(test.kinds) {
				t.Fatalf("Lex(%q) produced %d tokens (%v), want %d", test.src, len(tokens), tokens, len(test.kinds))
			}

			for i, tok := range tokens {
				if string(tok.Lexeme) != test.want[i] {
					t.Errorf("token %d: %q, want %q", i, tok.Lexeme, test.want[i])
				}

				if tok.Kind != test.kinds[i] {
					t.Errorf("token %d (%q): kind %v, want %v", i, tok.Lexeme, tok.Kind, test.kinds[i])
				}
			}
		})
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"unterminated comment", "1 (* oops", "unterminated comment"},
		{"stray character", "1 ; 2", "unrecognized character"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Lex(test.src)
			if err == nil {
				t.Fatalf("Lex(%q) succeeded, want error", test.src)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("Lex(%q) error = %v, want it to contain %q", test.src, err, test.want)
			}
		})
	}
}
