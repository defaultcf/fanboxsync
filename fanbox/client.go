package fanbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

	headers := []struct {
		key   string
		value string
	}{
		{
			key:   "Origin",
			value: fmt.Sprintf("https://%s.fanbox.cc", c.creatorId),
		},
		{
			key:   "Cookie",
			value: fmt.Sprintf("FANBOXSESSID=%s", c.sessionId),
		},
		{
			key:   "X-CSRF-Token",
			value: c.csrfToken,
		},
		{
			key:   "Content-Type",
			value: "application/json",
		},
	}
	for _, header := range headers {
		request.Header.Set(header.key, header.value)
	}

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

func (c *Client) PushPost(post *Post) error {
	url := "https://api.fanbox.cc/post.update"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	bodyJson, err := ConvertJson(post)
	if err != nil {
		return err
	}

	fields := []struct {
		key   string
		value string
	}{
		{
			key:   "postId",
			value: post.Id,
		},
		{
			key:   "status",
			value: string(post.Status),
		},
		{
			key:   "feeRequired",
			value: "500",
		},
		{
			key:   "title",
			value: post.Title,
		},
		{
			key:   "body",
			value: bodyJson,
		},
		{
			key:   "tags",
			value: "[]",
		},
		{
			key:   "tt",
			value: c.csrfToken,
		},
	}
	for _, field := range fields {
		err := writer.WriteField(field.key, field.value)
		if err != nil {
			return err
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return err
	}

	headers := []struct {
		key   string
		value string
	}{
		{
			key:   "Content-Type",
			value: writer.FormDataContentType(),
		},
		{
			key:   "Cookie",
			value: fmt.Sprintf("FANBOXSESSID=%s", c.sessionId),
		},
		{
			key:   "Origin",
			value: "https://www.fanbox.cc",
		},
	}
	for _, header := range headers {
		request.Header.Set(header.key, header.value)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("response: %+v", string(bytes))

	return nil
}
