package fanboxconv

import (
	"slices"
	"strings"
)

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
