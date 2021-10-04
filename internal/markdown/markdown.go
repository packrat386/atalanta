package markdown

import (
	"bytes"
)

func GenerateHTML(input []byte) []byte {
	return generateHTML(parseBlocks(sanitizeNewlines(input)))
}

func sanitizeNewlines(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}
