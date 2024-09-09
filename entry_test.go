package main_test

import (
	"testing"

	fanboxgo "github.com/defaultcf/fanbox-go"
	. "github.com/defaultcf/fanboxsync"
	"github.com/stretchr/testify/assert"
)

func TestConvertPost(t *testing.T) {
	tests := []struct {
		name string
		post fanboxgo.Post
		want Entry
	}{
		{
			name: "FANBOX から Markdown に変換できる",
			post: fanboxgo.Post{
				ID:          fanboxgo.NewOptString("1000000"),
				Title:       fanboxgo.NewOptString("テスト投稿"),
				Status:      fanboxgo.NewOptPostStatus(fanboxgo.PostStatusDraft),
				FeeRequired: fanboxgo.NewOptInt(500),
				Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
					Blocks: []fanboxgo.PostBodyBlocksItem{
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
							Text: fanboxgo.NewOptString("テキスト"),
						},
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeHeader),
							Text: fanboxgo.NewOptString("タイトル"),
						},
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
							Text: fanboxgo.NewOptString("これは太字です"),
							Styles: []fanboxgo.PostBodyBlocksItemStylesItem{
								{
									Type:   fanboxgo.NewOptPostBodyBlocksItemStylesItemType(fanboxgo.PostBodyBlocksItemStylesItemTypeBold),
									Offset: fanboxgo.NewOptInt(3),
									Length: fanboxgo.NewOptInt(2),
								},
							},
						},
					},
				}),
			},
			want: Entry{
				ID:     "1000000",
				Title:  "テスト投稿",
				Status: fanboxgo.PostStatusDraft,
				Fee:    "500",
				Body:   "テキスト\n## タイトル\nこれは**太字**です",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			e := Entry{}

			// execute
			e = *e.ConvertPost(&tt.post)

			// verify
			assert.Equal(t, tt.want, e)
		})
	}
}

func TestConvertFanbox(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
		want  fanboxgo.Post
	}{
		{
			name: "Markdown から FANBOX に変換できる",
			entry: Entry{
				ID:     "1000000",
				Title:  "テスト投稿",
				Status: fanboxgo.PostStatusDraft,
				Fee:    "500",
				Body:   "## ほげほげ\nテキスト\nこれは**太字**です",
			},
			want: fanboxgo.Post{
				ID:          fanboxgo.NewOptString("1000000"),
				Title:       fanboxgo.NewOptString("テスト投稿"),
				Status:      fanboxgo.NewOptPostStatus(fanboxgo.PostStatusDraft),
				FeeRequired: fanboxgo.NewOptInt(500),
				Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
					Blocks: []fanboxgo.PostBodyBlocksItem{
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeHeader),
							Text: fanboxgo.NewOptString("ほげほげ"),
						},
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
							Text: fanboxgo.NewOptString("テキスト"),
						},
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
							Text: fanboxgo.NewOptString("これは太字です"),
							Styles: []fanboxgo.PostBodyBlocksItemStylesItem{
								{
									Type:   fanboxgo.NewOptPostBodyBlocksItemStylesItemType(fanboxgo.PostBodyBlocksItemStylesItemTypeBold),
									Offset: fanboxgo.NewOptInt(3),
									Length: fanboxgo.NewOptInt(2),
								},
							},
						},
					},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			e := Entry{}

			// execute
			post := *e.ConvertFanbox(&tt.entry)

			// verify
			assert.Equal(t, tt.want, post)
		})
	}
}
