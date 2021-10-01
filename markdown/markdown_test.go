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
	ast := parse(mustReadFixture("test_doc.md"))

	for _, b := range ast.blocks {
		fmt.Println(b.kind)
		fmt.Printf("%#v\n", string(b.text))
	}
}
