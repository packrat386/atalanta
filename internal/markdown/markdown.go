package markdown

import (
	"bytes"
	"fmt"
)

func GenerateHTML(input []byte) (html []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	html = generateHTML(parseBlocks(sanitizeNewlines(input)))
	return
}

func sanitizeNewlines(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}
