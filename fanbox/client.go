package fanbox

import (
	"fmt"
	"net/http"
)

type Client struct {
	http.Client
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetPosts(creator_id string, session_id string) ([]*Post, error) {
	url := fmt.Sprintf("https://api.fanbox.cc/post.listCreator?creatorId=%s&limit=100", creator_id)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Origin", fmt.Sprintf("https://%s.fanbox.cc", creator_id))
	//request.Header.Set("Cookie", fmt.Sprintf("FANBOXSESSID=%s", session_id))

	response, err := c.Client.Do(request)
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
