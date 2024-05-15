package fanbox

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type fakeClient struct {
	http_client HttpClient
}

func NewFakeClient() *fakeClient {
	return &fakeClient{
		http_client: &fakeClient{},
	}
}

func (f *fakeClient) Do(request *http.Request) (*http.Response, error) {
	log.Print("This is Fake Do!")
	log.Printf("path: %s", request.URL.Path)

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
	log.Printf("response body: %+v", string(responseBody))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(string(responseBody))),
	}, nil
}
