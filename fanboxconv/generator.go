package fanboxconv

import (
	"fmt"
	"slices"
	"strings"
)

func isAllElmParentRoot(tokens []Token) bool {
	for _, token := range tokens {
		if token.ElmType != ElmTypeRoot {
			return false
		}
	}
	return true
}

// 閉じるタグの直後の index を返す...？
func getInsertPosition(content string) int {
	state := 0
	closeTagParentheses := []string{"<", ">"}
	position := 0

	for i, c := range content {
		if state == 1 && string(c) == closeTagParentheses[state] {
			position = i
			continue
		} else if state == 0 && string(c) == closeTagParentheses[state] {
			state += 1
		}
	}

	return position + 1
}

func createMergedContent(currentToken Token, parentToken Token) string {
	content := ""
	switch parentToken.ElmType {
	case ElmTypeBold:
		// TODO: 前のトークンまでの文字数を見て、上手く FANBOX 形式で表現する
		content = fmt.Sprintf("<bold>%s</bold>", currentToken.Content)
	case ElmTypeMerged:
		position := getInsertPosition(parentToken.Content)
		content = parentToken.Content[:position] + currentToken.Content + parentToken.Content[position:]
	}
	return content
}

func generateFanboxString(tokens []Token) string {
	contents := []string{}
	for _, token := range tokens {
		contents = append(contents, token.Content)
	}
	slices.Reverse(contents)
	return strings.Join(contents, "")
}

func generateFanbox(asts [][]Token) string {
	newAsts := []string{}
	for _, ast := range asts {
		slices.Reverse(ast)
		newAsts = append(newAsts, generateFanboxString(ast))
	}
	return strings.Join(newAsts, "")
}
