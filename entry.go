package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/defaultcf/fanboxsync/iframely"
)

type entry struct {
	iframelyClient *iframely.IframelyClient
	id             string
	title          string
	status         fanbox.PostStatus
	body           string
}

func NewEntry(id string, title string, status string, body string) *entry {
	return &entry{
		iframelyClient: iframely.NewIframelyClient(&http.Client{}),
		id:             id,
		title:          title,
		status:         fanbox.PostStatus(status),
		body:           body,
	}
}

func (e *entry) ConvertPost(post *fanbox.Post) *entry {
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
			url, err := e.getEmbedUrl(urlType, post.Body.UrlEmbedMap[v.UrlEmbedId])
			if err != nil {
				log.Fatal(err)
			} else {
				body = append(body, fmt.Sprintf("[%s](%s)", v.UrlEmbedId, url))
			}
		}
	}

	return &entry{
		id:     post.Id,
		title:  post.Title,
		status: post.Status,
		body:   strings.Join(body, "\n"),
	}
}

func (e *entry) ConvertFanbox() *fanbox.Post {
	// TODO: Markdown から FANBOX の形式に変換する

	return &fanbox.Post{}
}

func (e *entry) getEmbedUrl(urlType fanbox.UrlType, data fanbox.UrlEmbed) (string, error) {
	node, err := html.Parse(strings.NewReader(data.Html))
	if err != nil {
		return "", err
	}
	var url string
	switch urlType {
	case fanbox.UrlTypeCard:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url, err = e.iframelyClient.GetRealUrl(attr[0].Val)
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
