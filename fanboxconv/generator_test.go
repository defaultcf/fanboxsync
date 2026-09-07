package fanboxconv_test

import (
	"testing"

	. "github.com/defaultcf/fanboxsync/fanboxconv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFanbox(t *testing.T) {
	tests := []struct {
		name  string
		input []Token
		want  string
	}{
		{
			name: "正常に生成できる",
			input: []Token{
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
			want: "こんにちは",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup

			// execute
			str := GenerateFanboxString(tt.input)

			// verify
			assert.Equal(t, tt.want, str)
		})
	}
}
