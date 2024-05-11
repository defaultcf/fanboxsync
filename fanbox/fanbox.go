package fanbox

import (
	"encoding/json"
	"io"
)

type Post struct {
	Id    string
	Title string
}

type BodyPosts struct {
	Body struct {
		Items []*Post
	}
}

func ParsePosts(data io.Reader) ([]*Post, error) {
	bodyPosts := &BodyPosts{}
	bytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, bodyPosts)
	if err != nil {
		return nil, err
	}

	return bodyPosts.Body.Items, nil
}
