package markdown

import (
	"bytes"
	"fmt"
	"html"
	"io"
)

func generateHTML(ast *ast) []byte {
	var buf bytes.Buffer

	for _, b := range ast.blocks {
		writeBlockToHTML(b, &buf)
	}

	return buf.Bytes()
}

func writeBlockToHTML(b *block, w io.Writer) {
	switch b.kind {
	case blockBlank:
		writeBlockBlankToHTML(b, w)
	case blockH1:
		writeBlockH1ToHTML(b, w)
	case blockH2:
		writeBlockH2ToHTML(b, w)
	case blockH3:
		writeBlockH3ToHTML(b, w)
	case blockH4:
		writeBlockH4ToHTML(b, w)
	case blockH5:
		writeBlockH5ToHTML(b, w)
	case blockH6:
		writeBlockH6ToHTML(b, w)
	case blockQuote:
		writeBlockQuoteToHTML(b, w)
	case blockCode:
		writeBlockCodeToHTML(b, w)
	case blockParagraph:
		writeBlockParagraphToHTML(b, w)
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(b.kind)))
	}
}

func writeBlockBlankToHTML(b *block, w io.Writer) {
	return
}

func writeBlockH1ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h1>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))
	fmt.Fprintf(w, "</h1>")
}

func writeBlockH2ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h2>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</h2>")
}

func writeBlockH3ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h3>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</h3>")
}

func writeBlockH4ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h4>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</h4>")
}

func writeBlockH5ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h5>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</h5>")
}

func writeBlockH6ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h6>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</h6>")
}

func writeBlockQuoteToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<blockquote><p>")
	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</p></blockquote>")
}

func writeBlockCodeToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<pre><code>")

	b.text = bytes.TrimSuffix(b.text, []byte("\n"))
	b.text = bytes.TrimSuffix(b.text, []byte("```"))
	b.text = bytes.TrimPrefix(b.text, []byte("```\n"))

	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</code></pre>")
}

func writeBlockParagraphToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<p>")

	fmt.Fprintf(w, html.EscapeString(string(b.text)))

	fmt.Fprintf(w, "</p>")
}
