package fanbox

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type fakeClient struct {
	httpClient HttpClient
}

func NewFakeClient() *fakeClient {
	return &fakeClient{
		httpClient: nil,
	}
}

func (f *fakeClient) Do(request *http.Request) (*http.Response, error) {
	path := request.URL.Path
	var responseBody []byte
	switch path {
	case "/post.listManaged":
		responseBody, _ = json.Marshal(&BodyPosts{
			Body: []*Post{
				{
					Id:     "123456",
					Title:  "はじめての投稿",
					Status: "published",
				},
			},
		})
	case "/post.getEditable":
		responseBody, _ = json.Marshal(&BodyPost{
			Body: &Post{
				Id:     "123456",
				Title:  "はじめての投稿",
				Status: "published",
			},
		})
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(string(responseBody))),
	}, nil
}
