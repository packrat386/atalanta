package markdown

import ()

type AST struct {
	blocks []*block
}

func parse(input []byte) *AST {
	ast := &AST{blocks: []*block{}}

	parseBlocks(ast, input)

	return ast
}
