package markdown

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	doc := `## header
and a paragraph

and another with a \#hashtag
`
	toks, err := lex([]byte(doc))
	if err != nil {
		t.Fatal("got error: ", err)
	}

	for _, tok := range toks {
		fmt.Println(tok)
	}
}
