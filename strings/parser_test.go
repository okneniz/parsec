package strings_test

import (
	"os"
	"path/filepath"
	"testing"

	stdstrings "strings"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"literal text", "hello", "hello"},
		{"trailing input is left in the buffer", "hello world", "hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := strings.Parse([]rune(test.input), strings.String("expected text", "hello"))
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", test.input, err)
			}

			if got != test.want {
				t.Errorf("Parse(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"mismatched text", "world", "expected text"},
		{"empty input", "", "expected text"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := strings.Parse([]rune(test.input), strings.String("expected text", "hello"))
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want error", test.input)
			}

			if !stdstrings.Contains(err.Error(), test.want) {
				t.Errorf("Parse(%q) error = %v, want it to contain %q", test.input, err, test.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		parse parsec.Combinator[rune, strings.Position, string]
		want  string
	}{
		{"literal text", "hello", strings.String("expected text", "hello"), "hello"},
		{"multibyte runes", "привет", strings.String("expected text", "привет"), "привет"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := strings.ParseString(test.input, test.parse)
			if err != nil {
				t.Fatalf("ParseString(%q) failed: %v", test.input, err)
			}

			if got != test.want {
				t.Errorf("ParseString(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    uint64
	}{
		{"digits", "42", 42},
		{"long number", "123456", 123456},
		{"trailing text is left in the buffer", "42abc", 42},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "input.txt")

			if err := os.WriteFile(path, []byte(test.content), 0o600); err != nil {
				t.Fatalf("WriteFile failed: %v", err)
			}

			got, err := strings.ParseFile(path, strings.Unsigned[uint64]())
			if err != nil {
				t.Fatalf("ParseFile failed: %v", err)
			}

			if got != test.want {
				t.Errorf("ParseFile of %q = %d, want %d", test.content, got, test.want)
			}
		})
	}
}

func TestParseFileError(t *testing.T) {
	_, err := strings.ParseFile(
		filepath.Join(t.TempDir(), "missing.txt"),
		strings.String("expected text", "hello"),
	)
	if err == nil {
		t.Fatal("ParseFile succeeded, want error")
	}
}
