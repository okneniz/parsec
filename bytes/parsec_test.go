package bytes_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		parse parsec.Combinator[byte, int, uint16]
		want  uint16
	}{
		{"big endian uint16", []byte{0x01, 0x02}, bytes.ReadAs[uint16](2, "expected uint16", binary.BigEndian), 258},
		{"little endian uint16", []byte{0x01, 0x02}, bytes.ReadAs[uint16](2, "expected uint16", binary.LittleEndian), 513},
		{"trailing bytes are left in the buffer", []byte{0x01, 0x02, 0x03}, bytes.ReadAs[uint16](2, "expected uint16", binary.BigEndian), 258},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := bytes.Parse(test.input, test.parse)
			if err != nil {
				t.Fatalf("Parse(%v) failed: %v", test.input, err)
			}

			if got != test.want {
				t.Errorf("Parse(%v) = %d, want %d", test.input, got, test.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{"input shorter than the number", []byte{0x01}, "expected uint16"},
		{"empty input", nil, "expected uint16"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := bytes.Parse(test.input, bytes.ReadAs[uint16](2, "expected uint16", binary.BigEndian))
			if err == nil {
				t.Fatalf("Parse(%v) succeeded, want error", test.input)
			}

			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("Parse(%v) error = %v, want it to contain %q", test.input, err, test.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"ascii text", "hi!", "hi!"},
		{"a string converts to bytes, not runes", "\xd0\xb6", "\xd0\xb6"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parse := bytes.SequenceOf("expected text", []byte(test.input)...)

			got, err := bytes.ParseString(test.input, parse)
			if err != nil {
				t.Fatalf("ParseString(%q) failed: %v", test.input, err)
			}

			if string(got) != test.want {
				t.Errorf("ParseString(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"magic bytes", "\xca\xfe"},
		{"text", "hi!"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "input.bin")

			if err := os.WriteFile(path, []byte(test.content), 0o600); err != nil {
				t.Fatalf("WriteFile failed: %v", err)
			}

			parse := bytes.SequenceOf("expected file content", []byte(test.content)...)

			got, err := bytes.ParseFile(path, parse)
			if err != nil {
				t.Fatalf("ParseFile failed: %v", err)
			}

			if string(got) != test.content {
				t.Errorf("ParseFile of %q = %q, want %q", test.content, got, test.content)
			}
		})
	}
}

func TestParseFileError(t *testing.T) {
	_, err := bytes.ParseFile(
		filepath.Join(t.TempDir(), "missing.bin"),
		bytes.SequenceOf("expected magic bytes", 0xCA, 0xFE),
	)
	if err == nil {
		t.Fatal("ParseFile succeeded, want error")
	}
}
