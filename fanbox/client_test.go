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
			httpClient := fanbox.NewFakeClient()
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123")

			// execute
			posts, err := fanboxClient.GetPosts()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, posts)
		})
	}
}
