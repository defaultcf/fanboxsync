package fanbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type fakeClient struct {
	httpClient HttpClient
}

func NewFakeHttpClient() *fakeClient {
	return &fakeClient{
		httpClient: nil,
	}
}

func (f *fakeClient) Do(request *http.Request) (*http.Response, error) {
	var responseBody []byte
	switch request.URL.Path {
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
	case "/post.create":
		responseBody = []byte("{\"body\":{\"postId\":\"1234567\"}}")

	case "/post.update":
		responseBody = []byte("")
	default:
		return nil, fmt.Errorf("no paths matched")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(string(responseBody))),
	}, nil
}
