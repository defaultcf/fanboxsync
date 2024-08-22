package fanbox_test

import (
	"testing"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/stretchr/testify/assert"
)

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name string
		want []fanbox.Post
	}{
		{
			name: "投稿一覧を取得できる",
			want: []fanbox.Post{
				{
					Id:     "1234567",
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
			httpClient := fanbox.NewFakeHttpClient(map[string]fanbox.Post{
				"1234567": {
					Id:     "1234567",
					Title:  "はじめての投稿",
					Status: "published",
				},
			})
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			posts, err := fanboxClient.GetPosts()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, posts)
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name string
		want fanbox.Post
	}{
		{
			name: "指定した投稿を取得できる",
			want: fanbox.Post{
				Id:     "1234567",
				Title:  "はじめての投稿",
				Status: "published",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient(map[string]fanbox.Post{
				"1234567": {
					Id:     "1234567",
					Title:  "はじめての投稿",
					Status: "published",
				},
				"2345678": {
					Id:     "2345678",
					Title:  "次の投稿",
					Status: "published",
				},
			})
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			post, err := fanboxClient.GetPost("1234567")

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, post)
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "ポスト作成できる",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient(map[string]fanbox.Post{})
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			id, err := fanboxClient.CreatePost()

			// verify
			assert.NoError(t, err)
			posts, _ := fanboxClient.GetPosts()
			assert.Equal(t, tt.want, len(posts))

			_, err = fanboxClient.GetPost(id)
			assert.NoError(t, err)
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
				Id:          "1234567",
				Title:       "テスト",
				Status:      fanbox.PostStatusDraft,
				FeeRequired: 0,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			httpClient := fanbox.NewFakeHttpClient(map[string]fanbox.Post{
				"1234567": {
					Id:          "1234567",
					Title:       "これは変更前のタイトル",
					FeeRequired: 500,
				},
			})
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123", "csrfToken_123")

			// execute
			err := fanboxClient.PushPost(&tt.post)

			// verify
			assert.NoError(t, err)
			post, _ := fanboxClient.GetPost("1234567")
			assert.Equal(t, tt.post, post)
		})
	}
}
