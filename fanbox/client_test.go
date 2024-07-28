package fanbox_test

import (
	"testing"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/stretchr/testify/assert"
)

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name string
		want []*fanbox.Post
	}{
		{
			name: "投稿一覧を取得できる",
			want: []*fanbox.Post{
				{
					Id:     "123456",
					Title:  "はじめての投稿",
					Status: "published",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient()
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			posts, err := fanboxClient.GetPosts()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, posts)
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ポスト作成できる",
			want: "1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient()
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			postId, err := fanboxClient.CreatePost()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, postId)
		})
	}
}

func TestPushPost(t *testing.T) {
	tests := []struct {
		name string
		post fanbox.Post
		want error
	}{
		{
			name: "ポストを更新できる",
			post: fanbox.Post{
				Id:     "1234567",
				Title:  "テスト",
				Status: fanbox.PostStatusDraft,
				Body: fanbox.PostBody{
					Blocks: []fanbox.PostBodyBlock{
						{
							Type: fanbox.BodyTypeP,
							Text: "hoge",
						},
					},
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient()
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			err := fanboxClient.PushPost(&tt.post)

			// verify
			assert.NoError(t, err)
		})
	}
}
