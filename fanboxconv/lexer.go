package fanboxconv

import "regexp"

// 字句解析器（Lexer）

type Token struct {
	ID      int
	Parent  *Token
	ElmType ElmType
	Content string
}

type ElmType string

const (
	ElmTypeRoot   ElmType = "root"
	ElmTypeMerged ElmType = "merged"
	ElmTypeText   ElmType = "text"
	ElmTypeBold   ElmType = "bold"
)

const BoldElmRegexp = `\*\*(.+?)\*\*`

func genTextElement(id int, text string, parent Token) Token {
	return Token{
		ID:      id,
		Parent:  &parent,
		ElmType: ElmTypeText,
		Content: text,
	}
}

func genBoldElement(id int, text string, parent Token) Token {
	return Token{
		ID:      id,
		Parent:  &parent,
		ElmType: ElmTypeBold,
		Content: text,
	}
}

func indexWithBoldRegexp(text string) []int {
	re := regexp.MustCompile(BoldElmRegexp)
	return re.FindStringIndex(text)
}

func matchWithBoldRegexp(text string) []string {
	re := regexp.MustCompile(BoldElmRegexp)
	return re.FindStringSubmatch(text)
}
