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
			name:  "変換できる",
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
					Content: "hello",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup

			// execute
			tokens := Parse("hoge**hello**")

			// verify
			assert.Equal(t, tt.want, tokens)
		})
	}
}
