package markdown

func GenerateHTML(input []byte) []byte {
	return generateHTML(parseBlocks(input))
}
