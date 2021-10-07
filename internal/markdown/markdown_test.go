package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func mustReadFixture(filename string) []byte {
	buf, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		panic(err)
	}

	return buf
}

func TestParseBlocks(t *testing.T) {
	ast := parseBlocks(mustReadFixture("test_doc.md"))

	for _, b := range ast.blocks {
		fmt.Println(b.kind)
		fmt.Printf("%#v\n", string(b.text))
	}
}

func TestGenerateHTML(t *testing.T) {
	input := mustReadFixture("test_doc.md")

	html, err := GenerateHTML(input)
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	fmt.Println(string(html))
}

func TestInlineTextToHTML(t *testing.T) {
	text := []byte(`yadda _yadda_ yadda *yadda yadda*.`)

	sl := parseSpans(text)

	for sp := sl.begin; sp != nil; sp = sp.next {
		fmt.Printf("SPAN: %s\n", string(sp.text))
	}
}
