package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/defaultcf/fanboxsync/iframely"
)

type Entry struct {
	Id     string
	Title  string
	Status fanbox.PostStatus
	Body   string
}

func NewEntry(id string, title string, status string, body string) *Entry {
	return &Entry{
		Id:     id,
		Title:  title,
		Status: fanbox.PostStatus(status),
		Body:   body,
	}
}

func (e *Entry) ConvertPost(post *fanbox.Post) {
	var body []string
	for _, v := range post.Body.Blocks {
		switch v.Type {
		case fanbox.BodyTypeP:
			body = append(body, v.Text)
		case fanbox.BodyTypeHeader:
			body = append(body, fmt.Sprintf("## %s", v.Text))
		case fanbox.BodyTypeImage:
			body = append(body, fmt.Sprintf("![%s](%s)", v.ImageId, post.Body.ImageMap[v.ImageId].OriginalUrl))
		case fanbox.BodyTypeUrlEmbed:
			urlType := post.Body.UrlEmbedMap[v.UrlEmbedId].Type
			url, err := getEmbedUrl(urlType, post.Body.UrlEmbedMap[v.UrlEmbedId])
			if err != nil {
				log.Fatal(err)
			} else {
				body = append(body, fmt.Sprintf("[%s](%s)", v.UrlEmbedId, url))
			}
		}
	}
	e.Id = post.Id
	e.Title = post.Title
	e.Status = post.Status
	e.Body = strings.Join(body, "\n")
}

func (e *Entry) ConvertFanbox() *fanbox.Post {
	// TODO: Markdown から FANBOX の形式に変換する

	return &fanbox.Post{}
}

func getEmbedUrl(urlType fanbox.UrlType, data fanbox.UrlEmbed) (string, error) {
	node, err := html.Parse(strings.NewReader(data.Html))
	if err != nil {
		return "", err
	}
	var url string
	switch urlType {
	case fanbox.UrlTypeCard:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url, err = iframely.GetRealUrl(attr[0].Val)
		if err != nil {
			return "", err
		}
	case fanbox.UrlTypeHtml:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url = attr[0].Val
	case fanbox.UrlTypePost:
		url = fmt.Sprintf("https://%s.fanbox.cc/posts/%s", data.PostInfo.CreatorId, data.PostInfo.Id)
	case fanbox.UrlTypeDefault:
		url = data.Url
	default:
		return "", errors.New("unexpected url type")
	}

	return url, nil
}
