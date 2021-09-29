package main

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	doc := `## header
and a paragraph

and another with a \#hashtag
`
	toks, err := tokenize([]byte(doc))
	if err != nil {
		t.Fatal("got error: ", err)
	}

	for _, tok := range toks {
		fmt.Println(tok)
	}
}
