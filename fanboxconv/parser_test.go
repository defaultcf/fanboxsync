package fanboxconv_test

import (
	"testing"

	. "github.com/defaultcf/fanboxsync/fanboxconv"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "シンプルな Markdown を変換できる",
			input: "**こんにちは**",
			want: []Token{
				{
					ID: 1,
					Parent: &Token{
						ID:      0,
						Parent:  &Token{},
						ElmType: ElmTypeRoot,
					},
					ElmType: ElmTypeBold,
				},
				{
					ID: 2,
					Parent: &Token{
						ID: 1,
						Parent: &Token{
							ID:      0,
							Parent:  &Token{},
							ElmType: ElmTypeRoot,
						},
						ElmType: ElmTypeBold,
					},
					ElmType: ElmTypeText,
					Content: "こんにちは",
				},
			},
		},
		{
			name:  "左にテキストが来ても変換できる",
			input: "やっほー**こんにちは**",
			want: []Token{
				{
					ID: 1,
					Parent: &Token{
						ID:      0,
						Parent:  &Token{},
						ElmType: ElmTypeRoot,
					},
					ElmType: ElmTypeText,
					Content: "やっほー",
				},
				{
					ID: 2,
					Parent: &Token{
						ID:      0,
						Parent:  &Token{},
						ElmType: ElmTypeRoot,
					},
					ElmType: ElmTypeBold,
				},
				{
					ID: 3,
					Parent: &Token{
						ID: 2,
						Parent: &Token{
							ID:      0,
							Parent:  &Token{},
							ElmType: ElmTypeRoot,
						},
						ElmType: ElmTypeBold,
					},
					ElmType: ElmTypeText,
					Content: "こんにちは",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup

			// execute
			tokens := Parse(tt.input)

			// verify
			assert.Equal(t, tt.want, tokens)
		})
	}
}
