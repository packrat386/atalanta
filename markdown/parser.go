package markdown

type astNodeKind int64

const (
	astNodeDocument astNodeKind = iota
	astNodeH1
	astNodeH2
	astNodeH3
	astNodeH4
	astNodeH5
	astNodeH6
	astNodeParagraph
	astNodeParagraphText
)

type astNode struct {
	kind      astNodeKind
	parent    *astNode
	children  []*astNode
	beginning []token
	ending    []token
}

type ast struct {
	top  *astNode
	curr *astNode
}

type tokenStack struct {
	pos  int
	toks []token
}

func (t *tokenStack) guess(tk ...tokenKind) bool {
	for i := range tk {
		if t.pos+i < len(t.toks) {
			return false
		}

		if t.toks[t.pos+i].kind != tk[i] {
			return false
		}
	}

	return true
}

func (t *tokenStack) read(n int) []token {
	if t.pos+n >= len(t.toks) {
		panic("read too far")
	}

	ret := t.toks[t.pos : t.pos+n]

	t.pos += n

	return ret
}

func (t *tokenStack) rewind(n int) {
	if (t.pos - n) < 0 {
		panic("rewound too far")
	}

	t.pos -= n
}

func attemptClose(ast ast, ts *tokenStack) bool {

}

func attemptPreempt(ast ast, ts *tokenStack) bool {

}

func attemptNew(ast ast, ts *tokenStack) bool {

}
