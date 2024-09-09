package fanboxconv

import "strings"

func ConvertMarkdownToAST(markdown string) [][]Token {
	mdArray := strings.Split(markdown, "\n")
	asts := [][]Token{}
	for _, mdRow := range mdArray {
		asts = append(asts, parse(mdRow))
	}

	return asts
}
