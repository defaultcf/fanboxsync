package fanboxconv

import "strings"

var rootToken = Token{
	ID:      0,
	Parent:  &Token{},
	ElmType: ElmTypeRoot,
	Content: "",
}

func tokenize(id int, originalText string, parent Token) []Token {
	processingText := originalText
	elms := []Token{}
	p := parent

	for len(processingText) > 0 {
		matches := matchWithBoldRegexp(processingText)
		id += 1

		if len(matches) == 0 { // どこにもマッチしない
			elm := genTextElement(id, processingText, p)
			elms = append(elms, elm)
			processingText = ""
		} else {
			// TODO: 行頭にテキストが来ている場合、別に処理する

			elm := genBoldElement(id, "", p)
			elms = append(elms, elm)
			p = elm

			processingText = strings.Replace(processingText, matches[0], "", 1)
			elms = append(elms, tokenize(id, matches[1], p)...)
			p = parent
		}
	}

	return elms
}

func parse(markdownRow string) []Token {
	return tokenize(0, markdownRow, rootToken)
}
