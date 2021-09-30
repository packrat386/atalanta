package markdown

import (
	"fmt"
	"regexp"
)

type tokenKind int64

const (
	tokOctothorpe tokenKind = iota
	tokNewline
	tokSpace
	tokText
	tokEOF
)

func (k tokenKind) String() string {
	switch k {
	case tokOctothorpe:
		return "OCTOTHORPE"
	case tokNewline:
		return "NEWLINE"
	case tokSpace:
		return "SPACE"
	case tokText:
		return "TEXT"
	case tokEOF:
		return "EOF"
	default:
		panic(fmt.Sprint("unexpected tokenKind: ", int64(k)))
	}
}

type token struct {
	kind  tokenKind
	value []byte
}

func (t token) String() string {
	return fmt.Sprintf("%s: %s", t.kind.String(), string(t.value))
}

type tokenRule struct {
	kind    tokenKind
	matcher *regexp.Regexp
}

func readTokenRule(t tokenRule, buf []byte) (token, bool) {
	val := t.matcher.Find(buf)
	if val == nil {
		return token{}, false
	}

	return token{
		kind:  t.kind,
		value: val,
	}, true
}

var ruleset = []tokenRule{
	tokenRule{
		kind:    tokOctothorpe,
		matcher: regexp.MustCompile(`^#`),
	},
	tokenRule{
		kind:    tokNewline,
		matcher: regexp.MustCompile(`^(\n|\r\n|\r)`),
	},
	tokenRule{
		kind:    tokSpace,
		matcher: regexp.MustCompile(`^[\t\f ]+`),
	},
	tokenRule{
		kind:    tokText,
		matcher: regexp.MustCompile(`^(\\#|[^#\r\n\t\f ])+`),
	},
}

func readTokenRuleset(tt []tokenRule, buf []byte) (token, bool) {
	for _, t := range tt {
		tok, ok := readTokenRule(t, buf)
		if ok {
			return tok, true
		}
	}

	return token{}, false
}

func lex(buf []byte) ([]token, error) {
	pos := 0
	toks := []token{}

	for pos < len(buf) {
		tok, ok := readTokenRuleset(ruleset, buf[pos:])
		if !ok {
			return nil, fmt.Errorf("Could not tokenize: no valid token at position %d", pos)
		}

		fmt.Println("got one: ", tok)

		toks = append(toks, tok)
		pos += len(tok.value)
	}

	toks = append(toks, token{kind: tokEOF, value: []byte("")})

	return toks, nil
}
