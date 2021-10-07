package markdown

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustReadFixture(filename string) []byte {
	buf, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		panic(err)
	}

	return buf
}

func lineDiff(expected, actual []byte) string {
	s := new(strings.Builder)

	expectedL := bytes.SplitAfter(expected, []byte("\n"))
	actualL := bytes.SplitAfter(actual, []byte("\n"))

	max := len(expectedL)
	if max < len(actualL) {
		max = len(actualL)
	}

	for i := 0; i < max; i++ {
		if i >= len(expectedL) {
			fmt.Fprintln(s, "line: ", i)
			fmt.Fprintln(s, "not present in expected")
			fmt.Fprintf(s, "actual: %#v\n", string(actualL[i]))
		} else if i >= len(actualL) {
			fmt.Fprintln(s, "line: ", i)
			fmt.Fprintf(s, "expected : %#v\n", string(expectedL[i]))
			fmt.Fprintln(s, "not present in B")
		} else if !bytes.Equal(expectedL[i], actualL[i]) {
			fmt.Fprintln(s, "line: ", i)
			fmt.Fprintf(s, "expected : %#v\n", string(expectedL[i]))
			fmt.Fprintf(s, "actual   : %#v\n", string(actualL[i]))
		}
	}

	return s.String()
}

func TestGenerateHTML(t *testing.T) {
	tt := []string{
		"comprehensive",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"header_leading_space",
		"header_no_space",
		"header_escaped",
		"header_no_lazy",
		"blockquote",
		"blockquote_lazy",
		"code_block",
		"code_block_no_close",
		"paragraph",
		"extra_blank_lines",
		"inline_em",
		"inline_strong",
		"inline_em_strong_nested",
		"inline_em_no_match",
		"inline_em_precedence",
		"inline_link",
		"inline_link_bad_url",
		"inline_link_inside_em",
		"inline_link_breaks_em",
	}

	for _, tc := range tt {
		t.Run(tc, func(t *testing.T) {
			out, err := GenerateHTML(mustReadFixture(tc + ".md"))
			if err != nil {
				t.Fatalf("test case: '%s'\nunexpected error: %s", tc, err.Error())
			}

			expected := mustReadFixture(tc + ".out")

			if !bytes.Equal(out, expected) {
				t.Fatalf("test case: '%s'\ndiff:\n%s\n", tc, lineDiff(expected, out))
			}
		})
	}
}
