package markdown

import (
	"bytes"
	"fmt"
	"regexp"
)

type ast struct {
	blocks []*block
}

type block struct {
	kind blockKind
	text []byte
}

type blockKind int64

const (
	blockBlank blockKind = iota
	blockH1
	blockH2
	blockH3
	blockH4
	blockH5
	blockH6
	blockQuote
	blockCode
	blockParagraph
)

func (b blockKind) String() string {
	switch b {
	case blockBlank:
		return "BLANK"
	case blockH1:
		return "H1"
	case blockH2:
		return "H2"
	case blockH3:
		return "H3"
	case blockH4:
		return "H4"
	case blockH5:
		return "H5"
	case blockH6:
		return "H6"
	case blockQuote:
		return "QUOTE"
	case blockCode:
		return "CODE"
	case blockParagraph:
		return "PARAGRAPH"
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(b)))
	}
}

func parseBlocks(input []byte) *ast {
	ast := &ast{blocks: []*block{}}

	lines := bytes.SplitAfter(input, []byte("\n"))

	var current *block

	for _, l := range lines {
		if current == nil {
			current = openFirstBlock(l)
			ast.blocks = append(ast.blocks, current)
			continue
		}

		if canCloseBlock(current, l) {
			current = openFirstBlock(l)
			ast.blocks = append(ast.blocks, current)
			continue
		}

		continueBlock(current, l)
	}

	return ast
}

func openFirstBlock(line []byte) *block {
	for k := blockBlank; k <= blockParagraph; k++ {
		if canOpenBlock(k, line) {
			return openBlock(k, line)
		}
	}

	panic("could not open any block")
}

func canOpenBlock(k blockKind, line []byte) bool {
	switch k {
	case blockBlank:
		return canOpenBlockBlank(line)
	case blockH1:
		return canOpenBlockH1(line)
	case blockH2:
		return canOpenBlockH2(line)
	case blockH3:
		return canOpenBlockH3(line)
	case blockH4:
		return canOpenBlockH4(line)
	case blockH5:
		return canOpenBlockH5(line)
	case blockH6:
		return canOpenBlockH6(line)
	case blockQuote:
		return canOpenBlockQuote(line)
	case blockCode:
		return canOpenBlockCode(line)
	case blockParagraph:
		return canOpenBlockParagraph(line)
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(k)))
	}
}

func openBlock(k blockKind, line []byte) *block {
	switch k {
	case blockBlank:
		return openBlockBlank(line)
	case blockH1:
		return openBlockH1(line)
	case blockH2:
		return openBlockH2(line)
	case blockH3:
		return openBlockH3(line)
	case blockH4:
		return openBlockH4(line)
	case blockH5:
		return openBlockH5(line)
	case blockH6:
		return openBlockH6(line)
	case blockQuote:
		return openBlockQuote(line)
	case blockCode:
		return openBlockCode(line)
	case blockParagraph:
		return openBlockParagraph(line)
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(k)))
	}
}

func canCloseBlock(b *block, line []byte) bool {
	switch b.kind {
	case blockBlank:
		return canCloseBlockBlank(b, line)
	case blockH1:
		return canCloseBlockH1(b, line)
	case blockH2:
		return canCloseBlockH2(b, line)
	case blockH3:
		return canCloseBlockH3(b, line)
	case blockH4:
		return canCloseBlockH4(b, line)
	case blockH5:
		return canCloseBlockH5(b, line)
	case blockH6:
		return canCloseBlockH6(b, line)
	case blockQuote:
		return canCloseBlockQuote(b, line)
	case blockCode:
		return canCloseBlockCode(b, line)
	case blockParagraph:
		return canCloseBlockParagraph(b, line)
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(b.kind)))
	}
}

func continueBlock(b *block, line []byte) {
	switch b.kind {
	case blockBlank:
		continueBlockBlank(b, line)
	case blockH1:
		continueBlockH1(b, line)
	case blockH2:
		continueBlockH2(b, line)
	case blockH3:
		continueBlockH3(b, line)
	case blockH4:
		continueBlockH4(b, line)
	case blockH5:
		continueBlockH5(b, line)
	case blockH6:
		continueBlockH6(b, line)
	case blockQuote:
		continueBlockQuote(b, line)
	case blockCode:
		continueBlockCode(b, line)
	case blockParagraph:
		continueBlockParagraph(b, line)
	default:
		panic(fmt.Sprint("unrecognized block kind: ", int64(b.kind)))
	}
}

func trimPrefix(buf []byte, prefix *regexp.Regexp) []byte {
	idx := prefix.FindIndex(buf)
	if idx == nil {
		return buf
	}

	return buf[idx[1]:]
}

var blankLine = regexp.MustCompile(`^\s*$`)

func lineIsBlank(line []byte) bool {
	return blankLine.Match(line)
}

func canOpenBlockBlank(line []byte) bool {
	return lineIsBlank(line)
}

func openBlockBlank(line []byte) *block {
	return &block{
		kind: blockBlank,
		text: line,
	}
}

func canCloseBlockBlank(b *block, line []byte) bool {
	return true
}

func continueBlockBlank(b *block, line []byte) {
	return
}

var blockH1Opening = regexp.MustCompile(`^#\s+`)

func canOpenBlockH1(line []byte) bool {
	return blockH1Opening.Match(line)
}

func openBlockH1(line []byte) *block {
	return &block{
		kind: blockH1,
		text: trimPrefix(line, blockH1Opening),
	}
}

func canCloseBlockH1(b *block, line []byte) bool {
	return true
}

func continueBlockH1(b *block, line []byte) {
	return
}

var blockH2Opening = regexp.MustCompile(`^##\s+`)

func canOpenBlockH2(line []byte) bool {
	return blockH2Opening.Match(line)
}

func openBlockH2(line []byte) *block {
	return &block{
		kind: blockH2,
		text: trimPrefix(line, blockH2Opening),
	}
}

func canCloseBlockH2(b *block, line []byte) bool {
	return true
}

func continueBlockH2(b *block, line []byte) {
	return
}

var blockH3Opening = regexp.MustCompile(`^###\s+`)

func canOpenBlockH3(line []byte) bool {
	return blockH3Opening.Match(line)
}

func openBlockH3(line []byte) *block {
	return &block{
		kind: blockH3,
		text: trimPrefix(line, blockH3Opening),
	}
}

func canCloseBlockH3(b *block, line []byte) bool {
	return true
}

func continueBlockH3(b *block, line []byte) {
	return
}

var blockH4Opening = regexp.MustCompile(`^####\s+`)

func canOpenBlockH4(line []byte) bool {
	return blockH4Opening.Match(line)
}

func openBlockH4(line []byte) *block {
	return &block{
		kind: blockH4,
		text: trimPrefix(line, blockH4Opening),
	}
}

func canCloseBlockH4(b *block, line []byte) bool {
	return true
}

func continueBlockH4(b *block, line []byte) {
	return
}

var blockH5Opening = regexp.MustCompile(`^#####\s+`)

func canOpenBlockH5(line []byte) bool {
	return blockH5Opening.Match(line)
}

func openBlockH5(line []byte) *block {
	return &block{
		kind: blockH5,
		text: trimPrefix(line, blockH5Opening),
	}
}

func canCloseBlockH5(b *block, line []byte) bool {
	return true
}

func continueBlockH5(b *block, line []byte) {
	return
}

var blockH6Opening = regexp.MustCompile(`^######\s+`)

func canOpenBlockH6(line []byte) bool {
	return blockH6Opening.Match(line)
}

func openBlockH6(line []byte) *block {
	return &block{
		kind: blockH6,
		text: trimPrefix(line, blockH6Opening),
	}
}

func canCloseBlockH6(b *block, line []byte) bool {
	return true
}

func continueBlockH6(b *block, line []byte) {
	return
}

func canOpenBlockCode(line []byte) bool {
	return bytes.Equal(line, []byte("```\n"))
}

func openBlockCode(line []byte) *block {
	return &block{
		kind: blockCode,
		text: line,
	}
}

func canCloseBlockCode(b *block, line []byte) bool {
	return bytes.HasSuffix(b.text, []byte("\n```\n"))
}

func continueBlockCode(b *block, line []byte) {
	b.text = append(b.text, line...)
}

var blockQuoteOpening = regexp.MustCompile(`^>\s+`)

func canOpenBlockQuote(line []byte) bool {
	return blockQuoteOpening.Match(line)
}

func openBlockQuote(line []byte) *block {
	return &block{
		kind: blockQuote,
		text: trimPrefix(line, blockQuoteOpening),
	}
}

func canCloseBlockQuote(b *block, line []byte) bool {
	return canOpenBlockBlank(line) ||
		canOpenBlockH1(line) ||
		canOpenBlockH2(line) ||
		canOpenBlockH3(line) ||
		canOpenBlockH4(line) ||
		canOpenBlockH5(line) ||
		canOpenBlockH6(line) ||
		canOpenBlockCode(line)
}

func continueBlockQuote(b *block, line []byte) {
	b.text = append(b.text, trimPrefix(line, blockQuoteOpening)...)
}

func canOpenBlockParagraph(line []byte) bool {
	return true
}

func openBlockParagraph(line []byte) *block {
	return &block{
		kind: blockParagraph,
		text: line,
	}
}

func canCloseBlockParagraph(b *block, line []byte) bool {
	return canOpenBlockBlank(line) ||
		canOpenBlockH1(line) ||
		canOpenBlockH2(line) ||
		canOpenBlockH3(line) ||
		canOpenBlockH4(line) ||
		canOpenBlockH5(line) ||
		canOpenBlockH6(line) ||
		canOpenBlockCode(line) ||
		canOpenBlockQuote(line)
}

func continueBlockParagraph(b *block, line []byte) {
	b.text = append(b.text, line...)
}
