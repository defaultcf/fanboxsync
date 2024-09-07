package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
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
	for _, block := range post.Body.Value.Blocks {
		runeText := []rune(block.Text.Value)
		processedText := ""
		strPointer := 0
		switch t, _ := block.Type.Get(); t {
		case fanbox.PostBodyBlocksItemTypeP:
			// offset で昇順ソート
			sort.SliceStable(block.Styles, func(i, j int) bool { return block.Styles[i].Offset.Value < block.Styles[j].Offset.Value })
			for _, style := range block.Styles {
				switch style.Type.Value {
				case "bold": // 現在のところ bold だけ確認されている
					processedText += string(runeText[strPointer:style.Offset.Value])
					nextPointer := style.Offset.Value + style.Length.Value
					// style.Offset.Value から + style.Length までを ** で囲む
					processedText += fmt.Sprintf("**%s**", string(runeText[style.Offset.Value:nextPointer]))
					strPointer = nextPointer
				default:
					log.Fatal("unknown style type")
				}
			}
			// 残りの部分を追加
			processedText += string(runeText[strPointer:])
			body = append(body, processedText)
		case fanbox.PostBodyBlocksItemTypeHeader:
			body = append(body, fmt.Sprintf("## %s", block.Text.Value))
		case fanbox.PostBodyBlocksItemTypeImage:
			body = append(body, fmt.Sprintf("![%s](%s)", block.ImageId.Value, post.Body.Value.ImageMap.Value[block.ImageId.Value].OriginalUrl.Value))
		case fanbox.PostBodyBlocksItemTypeURLEmbed:
			urlType := post.Body.Value.UrlEmbedMap.Value[block.UrlEmbedId.Value].Type.Value
			url, err := e.getEmbedUrl(urlType, post.Body.Value.UrlEmbedMap.Value[block.UrlEmbedId.Value])
			if err != nil {
				log.Fatal(err)
			} else {
				body = append(body, fmt.Sprintf("[%s](%s)", block.UrlEmbedId.Value, url))
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
		re = regexp.MustCompile(`\*\*.+?\*\*`)
		matchIndexes := re.FindAllStringIndex(v, -1)
		styles := []fanbox.PostBodyBlocksItemStylesItem{}
		for _, matchIndex := range matchIndexes {
			offset := len([]rune(v[:matchIndex[0]]))
			length := len([]rune(v[matchIndex[0]:matchIndex[1]]))
			styles = append(styles, fanbox.PostBodyBlocksItemStylesItem{
				Type:   fanbox.NewOptString("bold"),
				Offset: fanbox.NewOptInt(offset),
				Length: fanbox.NewOptInt(length),
			})
		}
		// styles が空ならそもそも付けて送ってはならないため
		if len(styles) > 0 {
			blocks = append(blocks, fanbox.PostBodyBlocksItem{
				Type:   fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeP),
				Text:   fanbox.NewOptString(v),
				Styles: styles,
			})
		} else {
			blocks = append(blocks, fanbox.PostBodyBlocksItem{
				Type: fanbox.NewOptPostBodyBlocksItemType(fanbox.PostBodyBlocksItemTypeP),
				Text: fanbox.NewOptString(v),
			})
		}
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
