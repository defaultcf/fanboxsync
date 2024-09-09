package fanbox_test

import (
	"testing"

	fanboxgo "github.com/defaultcf/fanbox-go"
	. "github.com/defaultcf/fanboxsync/fanbox"
	"github.com/stretchr/testify/assert"
)

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name  string
		posts map[string]fanboxgo.Post
		want  []fanboxgo.Post
	}{
		{
			name: "一覧が取得できる",
			posts: map[string]fanboxgo.Post{
				"1000000": {
					ID:    fanboxgo.NewOptString("1000000"),
					Title: fanboxgo.NewOptString("最初の投稿"),
				},
				"1000001": {
					ID:    fanboxgo.NewOptString("1000001"),
					Title: fanboxgo.NewOptString("2番目の投稿"),
				},
			},
			want: []fanboxgo.Post{
				{
					ID:    fanboxgo.NewOptString("1000000"),
					Title: fanboxgo.NewOptString("最初の投稿"),
				},
				{
					ID:    fanboxgo.NewOptString("1000001"),
					Title: fanboxgo.NewOptString("2番目の投稿"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			client := NewFakeFanbox(tt.posts)
			testFanbox := NewTestFanbox(client)

			// execute
			posts, err := testFanbox.GetPosts()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, posts)
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name  string
		posts map[string]fanboxgo.Post
		id    string
		want  fanboxgo.Post
	}{
		{
			name: "投稿を取得できる",
			posts: map[string]fanboxgo.Post{
				"1000000": {
					ID:    fanboxgo.NewOptString("1000000"),
					Title: fanboxgo.NewOptString("最初の投稿"),
				},
				"1000001": {
					ID:    fanboxgo.NewOptString("1000001"),
					Title: fanboxgo.NewOptString("2番目の投稿"),
				},
			},
			id: "1000000",
			want: fanboxgo.Post{
				ID:    fanboxgo.NewOptString("1000000"),
				Title: fanboxgo.NewOptString("最初の投稿"),
			},
		},
		{
			name:  "投稿が無ければエラーが返る",
			posts: map[string]fanboxgo.Post{},
			id:    "100000",
			want:  fanboxgo.Post{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			client := NewFakeFanbox(tt.posts)
			testFanbox := NewTestFanbox(client)

			// execute
			post, err := testFanbox.GetPost(tt.id)

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, post)
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name  string
		posts map[string]fanboxgo.Post
		want  struct {
			id  string
			num int
		}
	}{
		{
			name:  "投稿が作成される",
			posts: map[string]fanboxgo.Post{},
			want: struct {
				id  string
				num int
			}{
				id:  "1000000",
				num: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			client := NewFakeFanbox(tt.posts)
			testFanbox := NewTestFanbox(client)

			// execute
			res, err := testFanbox.CreatePost()

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want.id, res)
			posts, _ := testFanbox.GetPosts()
			assert.Equal(t, tt.want.num, len(posts))
		})
	}
}

func TestPushPost(t *testing.T) {
	tests := []struct {
		name  string
		posts map[string]fanboxgo.Post
		post  fanboxgo.Post
		want  fanboxgo.Post
	}{
		{
			name: "投稿が更新される",
			posts: map[string]fanboxgo.Post{
				"1000000": {
					ID:          fanboxgo.NewOptString("1000000"),
					Title:       fanboxgo.NewOptString("変更前のタイトル"),
					FeeRequired: fanboxgo.NewOptInt(0),
					Status:      fanboxgo.NewOptPostStatus("draft"),
					Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
						Blocks: []fanboxgo.PostBodyBlocksItem{
							{
								Type: fanboxgo.NewOptPostBodyBlocksItemType("p"),
								Text: fanboxgo.NewOptString("変更前のテキスト"),
							},
						},
					}),
				},
			},
			post: fanboxgo.Post{
				ID:          fanboxgo.NewOptString("1000000"),
				Title:       fanboxgo.NewOptString("変更後のタイトル"),
				FeeRequired: fanboxgo.NewOptInt(500),
				Status:      fanboxgo.NewOptPostStatus("draft"),
				Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
					Blocks: []fanboxgo.PostBodyBlocksItem{
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType("p"),
							Text: fanboxgo.NewOptString("変更後のテキスト"),
						},
					},
				}),
			},
			want: fanboxgo.Post{
				ID:          fanboxgo.NewOptString("1000000"),
				Title:       fanboxgo.NewOptString("変更後のタイトル"),
				FeeRequired: fanboxgo.NewOptInt(500),
				Status:      fanboxgo.NewOptPostStatus("draft"),
				Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
					Blocks: []fanboxgo.PostBodyBlocksItem{
						{
							Type: fanboxgo.NewOptPostBodyBlocksItemType("p"),
							Text: fanboxgo.NewOptString("変更後のテキスト"),
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
			client := NewFakeFanbox(tt.posts)
			testFanbox := NewTestFanbox(client)

			// execute
			res, err := testFanbox.PushPost(&tt.post)

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, res)

			res, _ = testFanbox.GetPost(tt.post.ID.Value)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name  string
		posts map[string]fanboxgo.Post
		id    string
		want  int
	}{
		{
			name: "投稿を削除できる",
			posts: map[string]fanboxgo.Post{
				"1000000": {
					ID:    fanboxgo.NewOptString("1000000"),
					Title: fanboxgo.NewOptString("最初の投稿"),
				},
			},
			id:   "1000000",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			client := NewFakeFanbox(tt.posts)
			testFanbox := NewTestFanbox(client)

			// execute
			err := testFanbox.DeletePost(tt.id)

			// verify
			assert.NoError(t, err)
			posts, _ := testFanbox.GetPosts()
			assert.Equal(t, tt.want, len(posts))
		})
	}
}
