package fanbox

import (
	"encoding/json"
	"io"
)

type PostStatus string

const (
	PostStatusPublished PostStatus = "published"
	PostStatusDraft     PostStatus = "draft"
)

type BodyType string

const (
	BodyTypeP        BodyType = "p"
	BodyTypeHeader   BodyType = "header"
	BodyTypeImage    BodyType = "image"
	BodyTypeUrlEmbed BodyType = "url_embed"
)

type Post struct {
	Id     string
	Title  string
	Status PostStatus
	Body   struct {
		Blocks []struct {
			Type       BodyType // types: p, header, image, url_embed
			Text       string   // only p, header
			ImageId    string   // only image
			UrlEmbedId string   // only url_embed
			Styles     []struct {
				Type   string
				Offset int
				Length int
			}
		}
		ImageMap map[string]struct {
			Id           string
			Extension    string
			OriginalUrl  string
			ThumbnailUrl string
		}
		UrlEmbedMap map[string]struct {
			Id   string
			Type string
			Html string
		}
	}
}

type BodyPosts struct {
	Body []*Post
}

type BodyPost struct {
	Body *Post
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

	return bodyPosts.Body, nil
}

func ParsePost(data io.Reader) (*Post, error) {
	bodyPost := &BodyPost{}
	bytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, bodyPost)
	if err != nil {
		return nil, err
	}

	return bodyPost.Body, nil
}
