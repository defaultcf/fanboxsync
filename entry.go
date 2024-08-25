package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/defaultcf/fanbox-go"
	"github.com/defaultcf/fanboxsync/iframely"
	"golang.org/x/net/html"
)

type Entry struct {
	iframelyClient *iframely.IframelyClient
	id             string
	title          string
	status         fanbox.PostStatus
	fee            string
	body           string
	updatedAt      string
	publishedAt    string
}

func NewEntry(id string, title string, status string, fee string, body string) *Entry {
	return &Entry{
		iframelyClient: iframely.NewIframelyClient(&http.Client{}),
		id:             id,
		title:          title,
		status:         fanbox.PostStatus(status),
		fee:            fee,
		body:           body,
	}
}

// Fanbox から Markdown の形式に変換する
func (e *Entry) ConvertPost(post *fanbox.Post) *Entry {
	var body []string
	for _, v := range post.Body.Value.Blocks {
		switch t, _ := v.Type.Get(); t {
		case fanbox.PostBodyBlocksItemTypeP:
			body = append(body, v.Text.Value)
		case fanbox.PostBodyBlocksItemTypeHeader:
			body = append(body, fmt.Sprintf("## %s", v.Text.Value))
		case fanbox.PostBodyBlocksItemTypeImage:
			body = append(body, fmt.Sprintf("![%s](%s)", v.ImageId.Value, post.Body.Value.ImageMap.Value[v.ImageId.Value].OriginalUrl.Value))
		case fanbox.PostBodyBlocksItemTypeURLEmbed:
			urlType := post.Body.Value.UrlEmbedMap.Value[v.UrlEmbedId.Value].Type.Value
			url, err := e.getEmbedUrl(urlType, post.Body.Value.UrlEmbedMap.Value[v.UrlEmbedId.Value])
			if err != nil {
				log.Fatal(err)
			} else {
				body = append(body, fmt.Sprintf("[%s](%s)", v.UrlEmbedId.Value, url))
			}
		}
	}

	return &Entry{
		id:          post.ID.Value,
		title:       post.Title.Value,
		status:      post.Status.Value,
		fee:         fmt.Sprint(post.FeeRequired.Value),
		body:        strings.Join(body, "\n"),
		updatedAt:   post.UpdatedAt.Value,
		publishedAt: post.PublishedAt.Value,
	}
}

func (e *Entry) ConvertFanbox(entry *Entry) *fanbox.Post {
	blocks := []fanbox.PostBodyBlocksItem{}
	for _, v := range strings.Split(entry.body, "\n") {
		// Header
		re := regexp.MustCompile(`^## (.+)`)
		matches := re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanbox.PostBodyBlocksItem{
				Type: fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeHeader),
				Text: fanbox.NewOptString(matches[1]),
			})
			continue
		}
		// Image
		re = regexp.MustCompile(`^!\[(.+)\]\((.+)\)`)
		matches = re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanbox.PostBodyBlocksItem{
				Type:    fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeImage),
				ImageId: fanbox.NewOptString(matches[1]),
			})
			continue
		}
		// UrlEmbed
		re = regexp.MustCompile(`^\[(.+)\]\((.+)\)`)
		matches = re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanbox.PostBodyBlocksItem{
				Type:       fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeURLEmbed),
				UrlEmbedId: fanbox.NewOptString(matches[1]),
			})
			continue
		}
		// p
		if v == "" {
			continue
		}
		blocks = append(blocks, fanbox.PostBodyBlocksItem{
			Type: fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeP),
			Text: fanbox.NewOptString(v),
		})
	}

	fee, err := strconv.Atoi(entry.fee)
	if err != nil {
		return nil
	}

	return &fanbox.Post{
		ID:          fanbox.NewOptString(entry.id),
		Title:       fanbox.NewOptString(entry.title),
		Status:      fanbox.NewOptPostStatus(entry.status),
		FeeRequired: fanbox.NewOptInt(fee),
		Body: fanbox.NewOptPostBody(fanbox.PostBody{
			Blocks: blocks,
		}),
	}
}

func (e *Entry) getEmbedUrl(urlType fanbox.PostBodyUrlEmbedMapItemType, data fanbox.PostBodyUrlEmbedMapItem) (string, error) {
	node, err := html.Parse(strings.NewReader(data.HTML.Value))
	if err != nil {
		return "", err
	}
	var url string
	switch urlType {
	case fanbox.PostBodyUrlEmbedMapItemTypeHTMLCard:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url, err = e.iframelyClient.GetRealUrl(attr[0].Val)
		if err != nil {
			return "", err
		}
	case fanbox.PostBodyUrlEmbedMapItemTypeHTML:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url = attr[0].Val
	case fanbox.PostBodyUrlEmbedMapItemTypeFanboxPost:
		url = fmt.Sprintf("https://%s.fanbox.cc/posts/%s", data.PostInfo.Value.CreatorId.Value, data.PostInfo.Value.ID.Value)
	case fanbox.PostBodyUrlEmbedMapItemTypeDefault:
		url = data.URL.Value
	default:
		return "", errors.New("unexpected url type")
	}

	return url, nil
}
