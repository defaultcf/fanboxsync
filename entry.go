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

	fanboxgo "github.com/defaultcf/fanbox-go"
	"github.com/defaultcf/fanboxsync/iframely"
	"golang.org/x/net/html"
)

type Entry struct {
	iframelyClient *iframely.IframelyClient
	id             string
	title          string
	status         fanboxgo.PostStatus
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
		status:         fanboxgo.PostStatus(status),
		fee:            fee,
		body:           body,
	}
}

// Fanbox から Markdown の形式に変換する
func (e *Entry) ConvertPost(post *fanboxgo.Post) *Entry {
	var body []string
	for _, block := range post.Body.Value.Blocks {
		runeText := []rune(block.Text.Value)
		processedText := ""
		strPointer := 0
		switch t, _ := block.Type.Get(); t {
		case fanboxgo.PostBodyBlocksItemTypeP:
			// offset で昇順ソート
			sort.SliceStable(block.Styles, func(i, j int) bool { return block.Styles[i].Offset.Value < block.Styles[j].Offset.Value })
			for _, style := range block.Styles {
				switch style.Type.Value {
				case "bold": // 現在のところ bold だけ確認されている
					processedText += string(runeText[strPointer:style.Offset.Value])
					nextPointer := style.Offset.Value + style.Length.Value
					processedText += fmt.Sprintf("**%s**", string(runeText[style.Offset.Value:nextPointer]))
					strPointer = nextPointer
				default:
					log.Fatal("unknown style type")
				}
			}
			// 残りの部分を追加
			processedText += string(runeText[strPointer:])
			body = append(body, processedText)
		case fanboxgo.PostBodyBlocksItemTypeHeader:
			body = append(body, fmt.Sprintf("## %s", block.Text.Value))
		case fanboxgo.PostBodyBlocksItemTypeImage:
			body = append(body, fmt.Sprintf("![%s](%s)", block.ImageId.Value, post.Body.Value.ImageMap.Value[block.ImageId.Value].OriginalUrl.Value))
		case fanboxgo.PostBodyBlocksItemTypeURLEmbed:
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

func (e *Entry) ConvertFanbox(entry *Entry) *fanboxgo.Post {
	blocks := []fanboxgo.PostBodyBlocksItem{}
	for _, v := range strings.Split(entry.body, "\n") {
		// Header
		re := regexp.MustCompile(`^## (.+)`)
		matches := re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanboxgo.PostBodyBlocksItem{
				Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeHeader),
				Text: fanboxgo.NewOptString(matches[1]),
			})
			continue
		}
		// Image
		re = regexp.MustCompile(`^!\[(.+)\]\((.+)\)`)
		matches = re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanboxgo.PostBodyBlocksItem{
				Type:    fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeImage),
				ImageId: fanboxgo.NewOptString(matches[1]),
			})
			continue
		}
		// UrlEmbed
		re = regexp.MustCompile(`^\[(.+)\]\((.+)\)`)
		matches = re.FindStringSubmatch(v)
		if len(matches) > 0 {
			blocks = append(blocks, fanboxgo.PostBodyBlocksItem{
				Type:       fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeURLEmbed),
				UrlEmbedId: fanboxgo.NewOptString(matches[1]),
			})
			continue
		}
		// p
		re = regexp.MustCompile(`\*\*(.+?)\*\*`)
		matchIndexes := re.FindAllStringIndex(v, -1) // ここで得られる位置は rune ではなく string のもの
		styles := []fanboxgo.PostBodyBlocksItemStylesItem{}
		processedText := ""
		strPointer := 0
		for _, matchIndex := range matchIndexes {
			processedText += v[strPointer:matchIndex[0]]
			offset := len([]rune(processedText))
			matchStr := re.FindStringSubmatch(v[matchIndex[0]:matchIndex[1]])[1]
			processedText += matchStr
			length := len([]rune(processedText)) - offset

			styles = append(styles, fanboxgo.PostBodyBlocksItemStylesItem{
				Type:   fanboxgo.NewOptString("bold"),
				Offset: fanboxgo.NewOptInt(offset),
				Length: fanboxgo.NewOptInt(length),
			})
			strPointer = matchIndex[1]
		}
		// 残りの部分を追加
		processedText += v[strPointer:]

		// styles が空ならそもそも付けて送ってはならないため
		if len(styles) > 0 {
			blocks = append(blocks, fanboxgo.PostBodyBlocksItem{
				Type:   fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
				Text:   fanboxgo.NewOptString(processedText),
				Styles: styles,
			})
		} else {
			blocks = append(blocks, fanboxgo.PostBodyBlocksItem{
				Type: fanboxgo.NewOptPostBodyBlocksItemType(fanboxgo.PostBodyBlocksItemTypeP),
				Text: fanboxgo.NewOptString(v),
			})
		}
	}

	fee, err := strconv.Atoi(entry.fee)
	if err != nil {
		return nil
	}

	return &fanboxgo.Post{
		ID:          fanboxgo.NewOptString(entry.id),
		Title:       fanboxgo.NewOptString(entry.title),
		Status:      fanboxgo.NewOptPostStatus(entry.status),
		FeeRequired: fanboxgo.NewOptInt(fee),
		Body: fanboxgo.NewOptPostBody(fanboxgo.PostBody{
			Blocks: blocks,
		}),
	}
}

func (e *Entry) getEmbedUrl(urlType fanboxgo.PostBodyUrlEmbedMapItemType, data fanboxgo.PostBodyUrlEmbedMapItem) (string, error) {
	node, err := html.Parse(strings.NewReader(data.HTML.Value))
	if err != nil {
		return "", err
	}
	var url string
	switch urlType {
	case fanboxgo.PostBodyUrlEmbedMapItemTypeHTMLCard:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url, err = e.iframelyClient.GetRealUrl(attr[0].Val)
		if err != nil {
			return "", err
		}
	case fanboxgo.PostBodyUrlEmbedMapItemTypeHTML:
		attr := node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.Attr
		url = attr[0].Val
	case fanboxgo.PostBodyUrlEmbedMapItemTypeFanboxPost:
		url = fmt.Sprintf("https://%s.fanbox.cc/posts/%s", data.PostInfo.Value.CreatorId.Value, data.PostInfo.Value.ID.Value)
	case fanboxgo.PostBodyUrlEmbedMapItemTypeDefault:
		url = data.URL.Value
	default:
		return "", errors.New("unexpected url type")
	}

	return url, nil
}
