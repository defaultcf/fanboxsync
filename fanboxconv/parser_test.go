package fanboxconv_test

import (
	"fmt"
	"testing"

	. "github.com/defaultcf/fanboxsync/fanboxconv"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "変換できる",
			input: "**こんにちは**",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup

			// execute
			tokens := Parse("hoge**hello**")
			fmt.Printf("%+v", tokens)

			// verify
		})
	}
}
