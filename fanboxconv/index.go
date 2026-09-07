package fanboxconv

import "strings"

func ConvertMarkdownToAST(markdown string) string {
	mdArray := strings.Split(markdown, "\n")
	asts := [][]Token{}
	for _, mdRow := range mdArray {
		asts = append(asts, parse(mdRow))
	}
	fanboxString := generateFanbox(asts)

	return fanboxString
}
