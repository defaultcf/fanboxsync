package fanboxconv

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
		matches := indexWithBoldRegexp(processingText)
		id += 1

		if len(matches) == 0 { // どこにもマッチしない
			elm := genTextElement(id, processingText, p)
			elms = append(elms, elm)
			processingText = ""
		} else {
			// 行頭にテキストが来ている場合、別に処理する
			if matches[0] > 0 {
				elm := genTextElement(id, processingText[:matches[0]], p)
				elms = append(elms, elm)
				id += 1

			}

			elm := genBoldElement(id, "", p)
			elms = append(elms, elm)
			p = elm

			nextText := matchWithBoldRegexp(processingText[matches[0]:matches[1]])[1]
			elms = append(elms, tokenize(id, nextText, p)...)

			processingText = processingText[matches[1]:]
			p = parent
		}
	}

	return elms
}

func parse(markdownRow string) []Token {
	return tokenize(0, markdownRow, rootToken)
}
