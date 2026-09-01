package xml

import (
	"testing"

	"github.com/okneniz/parsec/strings"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "empty element",
			src:  `<book/>`,
			want: []string{"empty book"},
		},
		{
			name: "start and end",
			src:  `<book></book>`,
			want: []string{"start book", "end book"},
		},
		{
			name: "nested elements with text",
			src:  `<a><b>hello</b> world</a>`,
			want: []string{"start a", "start b", `text "hello"`, "end b", `text " world"`, "end a"},
		},
		{
			name: "attributes in both quotes",
			src:  `<point x='1' y="2"/>`,
			want: []string{`empty point x="1" y="2"`},
		},
		{
			name: "entity in attribute value",
			src:  `<a title="a&amp;b"/>`,
			want: []string{`empty a title="a&b"`},
		},
		{
			name: "comment",
			src:  `<a><!-- note: --></a>`,
			want: []string{"start a", `comment " note: "`, "end a"},
		},
		{
			name: "cdata keeps markup as text",
			src:  `<a><![CDATA[x < y]]></a>`,
			want: []string{"start a", `cdata "x < y"`, "end a"},
		},
		{
			name: "processing instruction",
			src:  `<?xml version="1.0"?><a/>`,
			want: []string{`procinst xml " version=\"1.0\""`, "empty a"},
		},
		{
			name: "entities in character data",
			src:  `<a>&amp;&#65;&#x42;&lt;</a>`,
			want: []string{"start a", `text "&AB<"`, "end a"},
		},
		{
			name: "whitespace between tags is character data",
			src:  "<a>\n  <b x=\"1\"/>\n</a>\n",
			want: []string{"start a", `text "\n  "`, `empty b x="1"`, `text "\n"`, "end a", `text "\n"`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			events, err := Parse(test.src)
			require.NoError(t, err)

			got := make([]string, 0, len(events))
			for _, ev := range events {
				got = append(got, ev.String())
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "mismatched end tag",
			src:  `<a></b>`,
			want: "unexpected end tag </b>",
		},
		{
			name: "unclosed element",
			src:  `<a><b></b>`,
			want: "unclosed element <a>",
		},
		{
			name: "element after root",
			src:  `<a/><b/>`,
			want: "content after root element",
		},
		{
			name: "text after root",
			src:  `<a/>junk`,
			want: "content after root element",
		},
		{
			name: "unterminated comment",
			src:  `<a><!-- x`,
			want: "unterminated comment",
		},
		{
			name: "unterminated cdata",
			src:  `<a><![CDATA[x`,
			want: "unterminated CDATA section",
		},
		{
			name: "unterminated attribute value",
			src:  `<a x="1>`,
			want: "unterminated attribute value",
		},
		{
			name: "unknown entity",
			src:  `<a>&bogus;</a>`,
			want: "unknown entity &bogus;",
		},
		{
			name: "malformed character reference",
			src:  `<a>&#xZZ;</a>`,
			want: "malformed character reference",
		},
		{
			name: "unterminated entity",
			src:  `<a>&amp`,
			want: "unterminated entity &amp",
		},
		{
			name: "stray bracket",
			src:  `<a> < </a>`,
			want: "expected element name",
		},
		{
			name: "digit cannot open a name",
			src:  `<1a/>`,
			want: "expected element name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := Parse(test.src)
			require.ErrorContains(t, err, test.want)
		})
	}
}

func TestEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     string
		want    []string
		wantErr string
	}{
		{
			name: "ends silently at the end of the document",
			src:  `<a/>`,
			want: []string{"empty a"},
		},
		{
			name:    "stops at the malformed event",
			src:     `<a>&bad;</a>`,
			want:    []string{"start a"},
			wantErr: "unknown entity &bad;",
		},
		{
			name: "no shape checks in the stream",
			src:  `<a></b>`,
			want: []string{"start a", "end b"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			seq, err := Events()(strings.Buffer([]rune(test.src)))
			require.NoError(t, err)

			var got []string
			var last error

			for ev, err := range seq {
				if err != nil {
					last = err

					break
				}

				got = append(got, ev.String())
			}

			require.Equal(t, test.want, got)

			if test.wantErr == "" {
				require.NoError(t, last)
			} else {
				require.ErrorContains(t, last, test.wantErr)
			}
		})
	}
}
