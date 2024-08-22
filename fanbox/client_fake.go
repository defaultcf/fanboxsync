package fanbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type fakeClient struct {
	httpClient HttpClient
	posts      map[string]Post
}

func NewFakeHttpClient(posts map[string]Post) *fakeClient {
	return &fakeClient{
		httpClient: nil,
		posts:      posts,
	}
}

func (f *fakeClient) Do(request *http.Request) (*http.Response, error) {
	var responseBody []byte
	switch request.URL.Path {
	case "/post.listManaged":
		responseBody, _ = json.Marshal(BodyPosts{
			Body: maps.Values(f.posts),
		})
	case "/post.getEditable":
		id := request.URL.Query().Get("postId")
		responseBody, _ = json.Marshal(BodyPost{
			Body: f.posts[id],
		})
	case "/post.create":
		id := fmt.Sprint(1000000 + len(f.posts))
		f.posts[id] = Post{Id: id}
		responseBody = []byte(fmt.Sprintf("{\"body\":{\"postId\":\"%s\"}}", id))
	case "/post.update":
		id := request.FormValue("postId")
		fee, err := strconv.Atoi(request.FormValue("feeRequired"))
		if err != nil {
			return nil, fmt.Errorf("cant parse fee")
		}
		post, _ := ParsePost(strings.NewReader(request.FormValue("body")))
		f.posts[id] = Post{
			Id:          id,
			Title:       request.FormValue("title"),
			Status:      PostStatus(request.FormValue("status")),
			FeeRequired: fee,
			Body:        post.Body,
		}
		responseBody = []byte("")
	default:
		return nil, fmt.Errorf("no paths matched")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil
}
