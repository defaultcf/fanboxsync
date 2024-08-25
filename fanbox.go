package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/defaultcf/fanbox-go"
)

type SecurityStore struct {
	rawCsrfToken string
	rawSessionId string
}

func (s SecurityStore) CsrfToken(ctx context.Context, operationName string) (fanbox.CsrfToken, error) {
	return fanbox.CsrfToken{
		APIKey: s.rawCsrfToken,
	}, nil
}

func (s SecurityStore) SessionId(ctx context.Context, operationName string) (fanbox.SessionId, error) {
	return fanbox.SessionId{
		APIKey: s.rawSessionId,
	}, nil
}

type CustomFanbox struct {
	Client        *fanbox.Client
	SecurityStore SecurityStore
}

func newFanbox(config *config) (*CustomFanbox, error) {
	s := SecurityStore{
		rawCsrfToken: config.Default.CsrfToken,
		rawSessionId: config.Default.SessionId,
	}
	c, err := fanbox.NewClient("https://api.fanbox.cc", s)
	if err != nil {
		return &CustomFanbox{}, err
	}

	return &CustomFanbox{
		Client:        c,
		SecurityStore: s,
	}, err
}

func (f *CustomFanbox) GetPosts() ([]fanbox.Post, error) {
	res, err := f.Client.ListManagedPosts(context.TODO(), fanbox.ListManagedPostsParams{Origin: "https://www.fanbox.cc"})
	if err != nil {
		return nil, err
	}
	return res.(*fanbox.List).Body, nil
}

func (f *CustomFanbox) GetPost(postId string) (fanbox.Post, error) {
	res, err := f.Client.GetEditablePost(context.TODO(), fanbox.GetEditablePostParams{Origin: "https://www.fanbox.cc", PostId: postId})
	if err != nil {
		return fanbox.Post{}, err
	}
	return res.(*fanbox.Get).Body.Value, nil
}

func (f *CustomFanbox) CreatePost() (fanbox.CreateBody, error) {
	res, err := f.Client.CreatePost(context.TODO(), fanbox.NewOptCreatePostReq(fanbox.CreatePostReq{Type: fanbox.CreatePostReqTypeArticle}), fanbox.CreatePostParams{Origin: "https://www.fanbox.cc"})
	if err != nil {
		return fanbox.CreateBody{}, err
	}

	switch r := res.(type) {
	case *fanbox.Create:
		return r.Body.Value, nil
	default:
		return fanbox.CreateBody{}, errors.New("error on create")
	}
}

func (f *CustomFanbox) PushPost(post *fanbox.Post) (fanbox.Post, error) {
	bodyJson, err := convertJson(&post.Body.Value.Blocks)
	if err != nil {
		return fanbox.Post{}, err
	}
	res, err := f.Client.UpdatePost(context.TODO(), fanbox.NewOptUpdatePostReq(fanbox.UpdatePostReq{
		PostId:      post.ID,
		Status:      fanbox.NewOptUpdatePostReqStatus(fanbox.UpdatePostReqStatus(post.Status.Value)),
		FeeRequired: fanbox.NewOptString(string(post.FeeRequired.Value)),
		Title:       post.Title,
		Body:        fanbox.NewOptString(bodyJson),
		Tags:        []string{},
		Tt:          fanbox.NewOptString(f.SecurityStore.rawCsrfToken),
	}), fanbox.UpdatePostParams{Origin: "https://www.fanbox.cc"})

	switch r := res.(type) {
	case *fanbox.Update:
		return r.Body.Value, nil
	default:
		return fanbox.Post{}, errors.New("error on update")
	}
}

func convertJson(body *[]fanbox.PostBodyBlocksItem) (string, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)

	// 本文が空なら、"null" ではなく空の配列を返す
	if jsonString == "null" {
		return "[]", nil
	}
	return jsonString, nil
}
