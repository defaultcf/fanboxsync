package fanbox

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"

	fanboxgo "github.com/defaultcf/fanbox-go"
)

type fakeFanbox struct {
	customFanbox CustomFanbox
	posts        map[string]fanboxgo.Post
}

func NewFakeFanbox(posts map[string]fanboxgo.Post) *fakeFanbox {
	client := &fakeFanbox{}
	return &fakeFanbox{
		customFanbox: CustomFanbox{
			Client:        client,
			SecurityStore: SecurityStore{},
		},
		posts: posts,
	}
}

func (f *fakeFanbox) CreatePost(ctx context.Context, request fanboxgo.OptCreatePostReq, params fanboxgo.CreatePostParams) (fanboxgo.CreatePostRes, error) {
	id := fmt.Sprint(1000000 + len(f.posts))
	f.posts[string(id)] = fanboxgo.Post{ID: fanboxgo.NewOptString(id)}

	//res := &fanboxgo.CreatePostBadRequestApplicationJSON{}
	//j, _ := json.Marshal(fanboxgo.CreateBody{PostId: fanboxgo.NewOptString(id)})
	//res.UnmarshalJSON(j)

	return &fanboxgo.Create{
		Body: fanboxgo.NewOptCreateBody(
			fanboxgo.CreateBody{
				PostId: fanboxgo.NewOptString(id),
			},
		),
	}, nil
}

func (f fakeFanbox) GetEditablePost(ctx context.Context, params fanboxgo.GetEditablePostParams) (fanboxgo.GetEditablePostRes, error) {
	//get := fanboxgo.Get{Body: fanboxgo.NewOptPost(f.posts[params.PostId])}
	//j, _ := json.Marshal(fanboxgo.Post{})
	//get.UnmarshalJSON(j)

	return &fanboxgo.Get{Body: fanboxgo.NewOptPost(f.posts[params.PostId])}, nil
}

func (f fakeFanbox) ListManagedPosts(ctx context.Context, params fanboxgo.ListManagedPostsParams) (fanboxgo.ListManagedPostsRes, error) {
	return &fanboxgo.List{Body: slices.Collect(maps.Values(f.posts))}, nil
}

func (f fakeFanbox) UpdatePost(ctx context.Context, request fanboxgo.OptUpdatePostReq, params fanboxgo.UpdatePostParams) (fanboxgo.UpdatePostRes, error) {
	fee, err := strconv.Atoi(request.Value.FeeRequired.Value)
	if err != nil {
		return nil, err
	}

	blocks := []fanboxgo.PostBodyBlocksItem{}
	err = json.Unmarshal([]byte(request.Value.Body.Value), &blocks)
	if err != nil {
		return nil, err
	}

	f.posts[request.Value.PostId.Value] = fanboxgo.Post{
		ID:          request.Value.PostId,
		Status:      fanboxgo.NewOptPostStatus(fanboxgo.PostStatus(request.Value.Status.Value)),
		FeeRequired: fanboxgo.NewOptInt(fee),
		Title:       fanboxgo.NewOptString(request.Value.Title.Value),
		Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
			Blocks: blocks,
		}),
	}
	return &fanboxgo.Update{Body: fanboxgo.NewOptPost(f.posts[request.Value.PostId.Value])}, nil
}
