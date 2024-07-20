package fanbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HttpClient
	creatorId  string
	sessionId  string
	csrfToken  string
}

func NewClient(httpClient HttpClient, creatorId string, sessionId string, csrfToken string) *Client {
	return &Client{
		httpClient: httpClient,
		creatorId:  creatorId,
		sessionId:  sessionId,
		csrfToken:  csrfToken,
	}
}

func (c *Client) GetPosts() ([]*Post, error) {
	url := "https://api.fanbox.cc/post.listManaged"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", c.creatorId))
	request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", c.sessionId))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ParsePosts(response.Body)
}

func (c *Client) GetPost(post_id string) (*Post, error) {
	url := fmt.Sprintf("https://api.fanbox.cc/post.getEditable?postId=%s", post_id)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", c.creatorId))
	request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", c.sessionId))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ParsePost(response.Body)
}

type CreatePostOptions struct {
	Type string `json:"type"`
}

type CreatePostBody struct {
	Body struct {
		PostId string
	}
}

func (c *Client) CreatePost() (string, error) {
	url := "https://api.fanbox.cc/post.create"
	options := &CreatePostOptions{
		Type: "article",
	}
	optionsJson, _ := json.Marshal(options)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(optionsJson))
	if err != nil {
		return "", err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", c.creatorId))
	request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", c.sessionId))
	request.Header.Set("X-CSRF-Token", c.csrfToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	createPostResponse := &CreatePostBody{}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, createPostResponse)
	if err != nil {
		return "", err
	}

	return createPostResponse.Body.PostId, nil
}

//func (c *Client) PushPost(post *Post) error {
//	url := fmt.Sprintf("https://api.fanbox.cc/post.create")
//
//	// TODO: Request.MultipartForm にデータを格納できるかも
//}
