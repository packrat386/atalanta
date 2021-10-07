package markdown

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/url"
	"regexp"
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
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h1>\n")
}

func writeBlockH2ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h2>")
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h2>\n")
}

func writeBlockH3ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h3>")
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h3>\n")
}

func writeBlockH4ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h4>")
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h4>\n")
}

func writeBlockH5ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h5>")
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h5>\n")
}

func writeBlockH6ToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<h6>")
	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))
	fmt.Fprintf(w, "</h6>\n")
}

func writeBlockQuoteToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<blockquote><p>")

	w.Write(bytes.TrimSuffix(parseSpans(b.text).bytes(), []byte("\n")))

	fmt.Fprintf(w, "</p></blockquote>\n")
}

func writeBlockCodeToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<pre><code>")

	b.text = bytes.TrimSuffix(b.text, []byte("\n"))
	b.text = bytes.TrimSuffix(b.text, []byte("```"))
	b.text = bytes.TrimPrefix(b.text, []byte("```\n"))

	w.Write(escapeHTMLBytes(bytes.TrimSuffix(b.text, []byte("\n"))))

	fmt.Fprintf(w, "</code></pre>\n")
}

func writeBlockParagraphToHTML(b *block, w io.Writer) {
	fmt.Fprintf(w, "<p>")

	w.Write(bytes.TrimSuffix(parseSpans(b.text).bytes(), []byte("\n")))

	fmt.Fprintf(w, "</p>\n")
}

type span struct {
	text []byte
	prev *span
	next *span
}

type spanList struct {
	begin *span
	end   *span
}

func (s *spanList) push(text []byte) *span {
	text = escapeHTMLBytes(text)

	if s.begin == nil && s.end == nil {
		sp := &span{
			text: text,
			prev: nil,
			next: nil,
		}

		s.begin = sp
		s.end = sp

		return sp
	}

	sp := &span{
		text: text,
		prev: s.end,
		next: nil,
	}

	s.end.next = sp
	s.end = sp

	return sp
}

func (s *spanList) pop() *span {
	if s.begin == nil && s.end == nil {
		panic("pop on empty list")
	}

	sp := s.end

	if s.begin == s.end {
		s.begin = nil
		s.end = nil
	} else {
		sp.prev.next = nil
		s.end = sp.prev
	}

	return sp
}

func (s *spanList) bytes() []byte {
	var buf bytes.Buffer

	for sp := s.begin; sp != nil; sp = sp.next {
		buf.Write(sp.text)
	}

	return buf.Bytes()
}

type delimiterKind int64

const (
	delimiterAsterisk delimiterKind = iota
	delimiterUnderscore
	delimiterOpenBracket
)

func (d delimiterKind) String() string {
	switch d {
	case delimiterAsterisk:
		return "ASTERISK"
	case delimiterUnderscore:
		return "UNDERSCORE"
	case delimiterOpenBracket:
		return "OPEN_BRACKET"
	default:
		panic(fmt.Sprint("unrecognized delimiter kind: ", int64(d)))
	}
}

type delimiter struct {
	kind     delimiterKind
	span     *span
	canOpen  bool
	canClose bool
	prev     *delimiter
	next     *delimiter
}

func (d delimiter) String() string {
	return fmt.Sprintf("kind: %s span: %p canOpen: %t canClose: %t prev: %p next: %p", d.kind, d.span, d.canOpen, d.canClose, d.prev, d.next)
}

type delimiterStack struct {
	begin *delimiter
	end   *delimiter
}

func (d *delimiterStack) push(delim *delimiter) {
	fmt.Println("push time")
	printStack(d)

	if d.begin == nil && d.end == nil {
		d.begin = delim
		d.end = delim

		delim.prev = nil
		delim.next = nil
		return
	}

	d.end.next = delim
	delim.prev = d.end
	delim.next = nil
	d.end = delim
}

func (d *delimiterStack) truncateAt(delim *delimiter) {
	d.end = delim
	delim.next = nil
}

func (d *delimiterStack) rm(delim *delimiter) {
	if delim == d.begin {
		d.begin = delim.next
	} else {
		delim.prev.next = delim.next
	}

	if delim == d.end {
		d.end = delim.prev
	} else {
		delim.next.prev = delim.prev
	}
}

func parseSpans(input []byte) *spanList {
	pos := 0
	buf := []byte{}
	sl := &spanList{
		begin: nil,
		end:   nil,
	}

	ds := &delimiterStack{
		begin: nil,
		end:   nil,
	}

	for pos < len(input) {
		str := []byte{input[pos]}

		switch {
		case bytes.Equal(str, []byte(`*`)):
			if len(buf) != 0 {
				sl.push(buf)
				buf = []byte{}
			}

			pos = consumeAsterisk(sl, ds, input, pos)
		case bytes.Equal(str, []byte(`_`)):
			if len(buf) != 0 {
				sl.push(buf)
				buf = []byte{}
			}

			pos = consumeUnderscore(sl, ds, input, pos)
		case bytes.Equal(str, []byte(`[`)):
			if len(buf) != 0 {
				sl.push(buf)
				buf = []byte{}
			}

			pos = consumeOpenBracket(sl, ds, input, pos)
		case bytes.Equal(str, []byte(`]`)):
			if len(buf) != 0 {
				sl.push(buf)
				buf = []byte{}
			}

			pos = consumeCloseBracket(sl, ds, input, pos)
		case bytes.Equal(str, []byte(`\`)):
			if len(buf) != 0 {
				sl.push(buf)
				buf = []byte{}
			}

			pos = consumeBackslash(sl, input, pos)
		default:
			buf = append(buf, input[pos])
			pos++
		}
	}

	if len(buf) != 0 {
		sl.push(buf)
	}

	processEmphasis(ds)

	return sl
}

var (
	openingAsteriskMatcher = regexp.MustCompile(`\*[^\s]`)
	closingAsteriskMatcher = regexp.MustCompile(`[^\s]\*`)
)

func consumeAsterisk(sl *spanList, ds *delimiterStack, input []byte, pos int) int {
	span := sl.push([]byte(`*`))

	canOpen := false
	canClose := false

	if pos > 0 {
		canClose = closingAsteriskMatcher.Match(input[pos-1 : pos+1])
		fmt.Printf("testing close: %#v -> %t\n", string(input[pos-1:pos+1]), canClose)
	}

	if pos+1 < len(input) {
		canOpen = openingAsteriskMatcher.Match(input[pos : pos+2])
	}

	if canOpen || canClose {
		fmt.Println("did we ever?")
		ds.push(&delimiter{
			kind:     delimiterAsterisk,
			span:     span,
			canOpen:  canOpen,
			canClose: canClose,
		})
	}

	return pos + 1
}

var (
	openingUnderscoreMatcher = regexp.MustCompile(`_[^\s]`)
	closingUnderscoreMatcher = regexp.MustCompile(`[^\s]_`)
)

func consumeUnderscore(sl *spanList, ds *delimiterStack, input []byte, pos int) int {
	span := sl.push([]byte(`_`))

	canOpen := false
	canClose := false

	if pos > 0 {
		canClose = closingUnderscoreMatcher.Match(input[pos-1 : pos+1])
	}

	if pos+1 < len(input) {
		canOpen = openingUnderscoreMatcher.Match(input[pos : pos+2])
	}

	if canOpen || canClose {
		ds.push(&delimiter{
			kind:     delimiterUnderscore,
			span:     span,
			canOpen:  canOpen,
			canClose: canClose,
		})
	}
	return pos + 1
}

func consumeOpenBracket(sl *spanList, ds *delimiterStack, input []byte, pos int) int {
	span := sl.push([]byte(`[`))

	ds.push(&delimiter{
		kind:     delimiterOpenBracket,
		span:     span,
		canOpen:  true,
		canClose: false,
	})

	return pos + 1
}

var linkURLMatcher = regexp.MustCompile(`\]\(([^\s\)]*)\)`)

func consumeCloseBracket(sl *spanList, ds *delimiterStack, input []byte, pos int) int {
	span := sl.push([]byte(`]`))

	matches := linkURLMatcher.FindSubmatch(input[pos:])
	if matches == nil {
		return pos + 1
	}

	u, err := url.Parse(string(matches[1]))
	if err != nil {
		return pos + 1
	}

	var opener *delimiter
	for d := ds.end; d != nil; d = d.prev {
		if d.kind == delimiterOpenBracket {
			opener = d
			ds.truncateAt(opener)
			break
		}
	}

	if opener == nil {
		return pos + 1
	}

	opener.span.text = []byte(fmt.Sprintf(`<a href="%s">`, u.String()))
	span.text = []byte(`</a>`)
	return pos + len(matches[0])
}

func consumeBackslash(sl *spanList, input []byte, pos int) int {
	if pos+1 >= len(input) {
		// we're at the end
		sl.push([]byte(`\`))
		return pos + 1
	}

	next := []byte{input[pos+1]}

	switch {
	case bytes.Equal(next, []byte(`*`)):
		sl.push([]byte(`*`))
		return pos + 2
	case bytes.Equal(next, []byte(`_`)):
		sl.push([]byte(`_`))
		return pos + 2
	case bytes.Equal(next, []byte(`[`)):
		sl.push([]byte(`[`))
		return pos + 2
	case bytes.Equal(next, []byte(`]`)):
		sl.push([]byte(`]`))
		return pos + 2
	case bytes.Equal(next, []byte(`\`)):
		sl.push([]byte(`\`))
		return pos + 2
	default: // nothing to escape
		sl.push([]byte(`\`))
		return pos + 1
	}
}

func processEmphasis(ds *delimiterStack) {
	fmt.Println("process emphasis")
	printStack(ds)

	for closer := nextClosingDelimiter(ds); closer != nil; closer = nextClosingDelimiter(ds) {
		fmt.Printf("found closer: %p\n", closer)

		opener := openingDelimiter(ds, closer)
		if opener == nil {
			ds.rm(closer)

			fmt.Println("no opener")
			printStack(ds)
			continue
		}

		fmt.Printf("found opener: %p\n", opener)

		if opener.kind == delimiterUnderscore {
			fmt.Printf("process underscore %p <-> %p\n", opener, closer)
			opener.span.text = []byte("<em>")
			closer.span.text = []byte("</em>")
		}

		if opener.kind == delimiterAsterisk {
			fmt.Printf("process asterisk %p <-> %p\n", opener, closer)
			opener.span.text = []byte("<strong>")
			closer.span.text = []byte("</strong>")
		}

		rmBetween(ds, opener, closer)

		fmt.Println("after process")
		printStack(ds)
	}

	fmt.Println("done w/ emphasis")
	printStack(ds)
}

func openingDelimiter(ds *delimiterStack, closer *delimiter) *delimiter {
	fmt.Println("looking for opener")
	for d := closer.prev; d != nil; d = d.prev {
		fmt.Printf("checking %p\n", d)
		if d.canOpen && d.kind == closer.kind {
			return d
		}
	}

	return nil
}

func nextClosingDelimiter(ds *delimiterStack) *delimiter {
	for d := ds.begin; d != nil; d = d.next {
		if d.canClose {
			return d
		}
	}

	return nil
}

func rmBetween(ds *delimiterStack, begin *delimiter, end *delimiter) {
	ptr := begin
	for ptr != end.next {
		next := ptr.next
		ds.rm(ptr)
		ptr = next
	}
}

func printStack(ds *delimiterStack) {
	fmt.Println("BEGIN STACK")
	for d := ds.begin; d != nil; d = d.next {
		fmt.Printf("%p -> %s\n", d, d.String())
	}
	fmt.Println("END STACK")
}

func escapeHTMLBytes(in []byte) []byte {
	return []byte(html.EscapeString(string(in)))
}
