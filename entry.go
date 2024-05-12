package main

import (
	"fmt"
	"strings"

	"github.com/defaultcf/fanboxsync/fanbox"
)

type Entry struct {
	Id     string
	Title  string
	Status fanbox.PostStatus
	Body   string // TODO: Markdown で持たせる
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
	// TODO: FANBOX から Markdown の形式に変換する
	var body []string
	for _, v := range post.Body.Blocks {
		switch v.Type {
		case fanbox.BodyTypeP:
			body = append(body, v.Text)
		case fanbox.BodyTypeHeader:
			body = append(body, fmt.Sprintf("## %s", v.Text))
		case fanbox.BodyTypeImage:
			body = append(body, fmt.Sprintf("![](%s)", post.Body.ImageMap[v.ImageId].OriginalUrl))
		case fanbox.BodyTypeUrlEmbed:
			body = append(body, post.Body.UrlEmbedMap[v.UrlEmbedId].Html)
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
