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
	defaultParams defaultParams
}

type defaultParams struct {
	origin    string
	userAgent string
}

func NewFanbox(csrfToken, sessionId, userAgent string) (*CustomFanbox, error) {
	s := SecurityStore{
		rawCsrfToken: csrfToken,
		rawSessionId: sessionId,
	}
	d := defaultParams{
		origin:    "https://www.fanbox.cc",
		userAgent: userAgent,
	}
	c, err := fanboxgo.NewClient("https://api.fanbox.cc", s)
	if err != nil {
		return &CustomFanbox{}, err
	}

	return &CustomFanbox{
		Client:        c,
		SecurityStore: s,
		defaultParams: d,
	}, nil
}

func NewTestFanbox(client fanboxgo.Invoker) *CustomFanbox {
	return &CustomFanbox{
		Client: client,
	}
}

func (f CustomFanbox) GetPosts() ([]fanboxgo.Post, error) {
	res, err := f.Client.ListManagedPosts(context.TODO(), fanboxgo.ListManagedPostsParams{
		Origin:    f.defaultParams.origin,
		UserAgent: f.defaultParams.userAgent,
	})
	if err != nil {
		return nil, err
	}
	return res.(*fanboxgo.List).Body, nil
}

func (f CustomFanbox) GetPost(postId string) (fanboxgo.Post, error) {
	res, err := f.Client.GetEditablePost(context.TODO(), fanboxgo.GetEditablePostParams{
		Origin:    f.defaultParams.origin,
		UserAgent: f.defaultParams.userAgent,
		PostId:    postId,
	})
	if err != nil {
		return fanboxgo.Post{}, err
	}
	return res.(*fanboxgo.Get).Body.Value, nil
}

func (f CustomFanbox) CreatePost() (string, error) {
	res, err := f.Client.CreatePost(context.TODO(),
		fanboxgo.NewOptCreatePostReq(fanboxgo.CreatePostReq{Type: fanboxgo.CreatePostReqTypeArticle}),
		fanboxgo.CreatePostParams{
			Origin:    f.defaultParams.origin,
			UserAgent: f.defaultParams.userAgent,
		},
	)
	if err != nil {
		return "", err
	}

	switch r := res.(type) {
	case *fanboxgo.Create:
		return r.Body.Value.PostId.Value, nil
	default:
		return "", errors.New("error on create")
	}
}

func (f CustomFanbox) PushPost(post *fanboxgo.Post) (fanboxgo.Post, error) {
	var commentingPermissionScope fanboxgo.UpdatePostReqCommentingPermissionScope
	if post.FeeRequired.Value == 0 {
		commentingPermissionScope = fanboxgo.UpdatePostReqCommentingPermissionScopeEveryone
	} else {
		commentingPermissionScope = fanboxgo.UpdatePostReqCommentingPermissionScopeSupporters
	}

	bodyJson, err := convertJson(&post.Body.Value.Blocks)
	if err != nil {
		return fanboxgo.Post{}, err
	}
	res, err := f.Client.UpdatePost(context.TODO(),
		fanboxgo.NewOptUpdatePostReq(fanboxgo.UpdatePostReq{
			PostId:                    post.ID,
			Status:                    fanboxgo.NewOptUpdatePostReqStatus(fanboxgo.UpdatePostReqStatus(post.Status.Value)),
			FeeRequired:               fanboxgo.NewOptString(fmt.Sprint(post.FeeRequired.Value)),
			Title:                     post.Title,
			CommentingPermissionScope: fanboxgo.NewOptUpdatePostReqCommentingPermissionScope(commentingPermissionScope),
			Body:                      fanboxgo.NewOptString(bodyJson),
			Tags:                      []string{},
			Tt:                        fanboxgo.NewOptString(f.SecurityStore.rawCsrfToken),
		}),
		fanboxgo.UpdatePostParams{
			Origin:    f.defaultParams.origin,
			UserAgent: f.defaultParams.userAgent,
		},
	)
	if err != nil {
		return fanboxgo.Post{}, err
	}

	switch r := res.(type) {
	case *fanboxgo.Update:
		return r.Body.Value, nil
	default:
		return fanboxgo.Post{}, errors.New("error on update")
	}
}

func (f CustomFanbox) DeletePost(postId string) error {
	_, err := f.Client.DeletePost(
		context.TODO(),
		fanboxgo.NewOptDeletePostReq(fanboxgo.DeletePostReq{PostId: postId}),
		fanboxgo.DeletePostParams{
			Origin:    f.defaultParams.origin,
			UserAgent: f.defaultParams.userAgent,
		},
	)
	if err != nil {
		return err
	}
	return nil
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
