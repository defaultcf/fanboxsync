package fanbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	fanboxgo "github.com/defaultcf/fanbox-go"
)

type SecurityStore struct {
	rawCsrfToken string
	rawSessionId string
}

func (s SecurityStore) CsrfToken(ctx context.Context, operationName string) (fanboxgo.CsrfToken, error) {
	return fanboxgo.CsrfToken{
		APIKey: s.rawCsrfToken,
	}, nil
}

func (s SecurityStore) SessionId(ctx context.Context, operationName string) (fanboxgo.SessionId, error) {
	return fanboxgo.SessionId{
		APIKey: s.rawSessionId,
	}, nil
}

type CustomFanbox struct {
	Client        fanboxgo.Invoker
	SecurityStore SecurityStore
}

func NewFanbox(csrfToken, sessionId string) (*CustomFanbox, error) {
	s := SecurityStore{
		rawCsrfToken: csrfToken,
		rawSessionId: sessionId,
	}
	c, err := fanboxgo.NewClient("https://api.fanbox.cc", s)
	if err != nil {
		return &CustomFanbox{}, err
	}

	return &CustomFanbox{
		Client:        c,
		SecurityStore: s,
	}, err
}

func (f *CustomFanbox) GetPosts() ([]fanboxgo.Post, error) {
	res, err := f.Client.ListManagedPosts(context.TODO(), fanboxgo.ListManagedPostsParams{Origin: "https://www.fanbox.cc"})
	if err != nil {
		return nil, err
	}
	return res.(*fanboxgo.List).Body, nil
}

func (f *CustomFanbox) GetPost(postId string) (fanboxgo.Post, error) {
	res, err := f.Client.GetEditablePost(context.TODO(), fanboxgo.GetEditablePostParams{Origin: "https://www.fanbox.cc", PostId: postId})
	if err != nil {
		return fanboxgo.Post{}, err
	}
	return res.(*fanboxgo.Get).Body.Value, nil
}

func (f *CustomFanbox) CreatePost() (fanboxgo.CreateBody, error) {
	res, err := f.Client.CreatePost(context.TODO(), fanboxgo.NewOptCreatePostReq(fanboxgo.CreatePostReq{Type: fanboxgo.CreatePostReqTypeArticle}), fanboxgo.CreatePostParams{Origin: "https://www.fanbox.cc"})
	if err != nil {
		return fanboxgo.CreateBody{}, err
	}

	switch r := res.(type) {
	case *fanboxgo.Create:
		return r.Body.Value, nil
	default:
		return fanboxgo.CreateBody{}, errors.New("error on create")
	}
}

func (f *CustomFanbox) PushPost(post *fanboxgo.Post) (fanboxgo.Post, error) {
	bodyJson, err := convertJson(&post.Body.Value.Blocks)
	if err != nil {
		return fanboxgo.Post{}, err
	}
	res, err := f.Client.UpdatePost(context.TODO(), fanboxgo.NewOptUpdatePostReq(fanboxgo.UpdatePostReq{
		PostId:      post.ID,
		Status:      fanboxgo.NewOptUpdatePostReqStatus(fanboxgo.UpdatePostReqStatus(post.Status.Value)),
		FeeRequired: fanboxgo.NewOptString(fmt.Sprint(post.FeeRequired.Value)),
		Title:       post.Title,
		Body:        fanboxgo.NewOptString(bodyJson),
		Tags:        []string{},
		Tt:          fanboxgo.NewOptString(f.SecurityStore.rawCsrfToken),
	}), fanboxgo.UpdatePostParams{Origin: "https://www.fanbox.cc"})

	switch r := res.(type) {
	case *fanboxgo.Update:
		return r.Body.Value, nil
	default:
		return fanboxgo.Post{}, errors.New("error on update")
	}
}

func convertJson(body *[]fanboxgo.PostBodyBlocksItem) (string, error) {
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
