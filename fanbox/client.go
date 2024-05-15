package fanbox

import (
	"fmt"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http_client HttpClient
	creator_id  string
	session_id  string
}

func NewClient(http_client HttpClient, creator_id string, session_id string) *Client {
	return &Client{
		http_client: http_client,
		creator_id:  creator_id,
		session_id:  session_id,
	}
}

func (c *Client) GetPosts() ([]*Post, error) {
	url := "https://api.fanbox.cc/post.listManaged"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", c.creator_id))
	request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", c.session_id))

	response, err := c.http_client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	posts, err := ParsePosts(response.Body)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *Client) GetPost(post_id string) (*Post, error) {
	url := fmt.Sprintf("https://api.fanbox.cc/post.getEditable?postId=%s", post_id)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", c.creator_id))
	request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", c.session_id))

	response, err := c.http_client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	post, err := ParsePost(response.Body)
	if err != nil {
		return nil, err
	}

	return post, nil
}
