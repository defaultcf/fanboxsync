package fanbox_test

import (
	"log"
	"testing"

	"github.com/defaultcf/fanboxsync/fanbox"
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
					Id:     "123456",
					Title:  "はじめての投稿",
					Status: "published",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpClient := fanbox.NewFakeClient()
			fanboxClient := fanbox.NewClient(httpClient, "creator_123", "session_123")

			posts, err := fanboxClient.GetPosts()
			log.Printf("posts: %+v", posts)
			if err != nil {
				t.Error("実行に失敗しました")
			}
			if posts[0].Id != test.want[0].Id {
				t.Error("投稿一覧が正しく取得できていません")
			}
		})
	}
}
