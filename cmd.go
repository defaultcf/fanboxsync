package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/goccy/go-yaml"
)

func CommandPull(config *config) error {
	client := fanbox.NewClient(
		&http.Client{},
		config.Default.CreatorId,
		config.Default.SessionId,
		config.Default.CsrfToken,
	)
	posts, err := client.GetPosts()
	if err != nil {
		return err
	}
	for _, v := range posts {
		post, err := client.GetPost(v.Id)
		if err != nil {
			return err
		}

		e := NewEntry("", "", "", "")
		convertedEntry := e.ConvertPost(post)

		err = saveFile(*convertedEntry)
		if err != nil {
			return err
		}
	}
	return nil
}

func CommandCreate(config *config, title string) error {
	client := fanbox.NewClient(
		&http.Client{},
		config.Default.CreatorId,
		config.Default.SessionId,
		config.Default.CsrfToken,
	)
	postId, err := client.CreatePost()
	if err != nil {
		return err
	}
	entry := NewEntry(postId, title, string(fanbox.PostStatusDraft), "")
	post := entry.ConvertFanbox(entry)
	client.PushPost(post) // タイトルをセット

	entry.updatedAt = time.Now().Format(time.RFC3339)
	err = saveFile(*entry)
	if err != nil {
		return err
	}

	return nil
}

type meta struct {
	Id     string `yaml:"id"`
	Title  string `yaml:"title"`
	Status string `yaml:"status"`
}

func CommandPush(config *config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// メタデータをマークダウンから抽出
	rawBody := string(bytes)
	reMeta := regexp.MustCompile(`---\n`)
	splited := reMeta.Split(rawBody, 3)
	m := meta{}
	err = yaml.Unmarshal([]byte(splited[1]), &m)
	if err != nil {
		return err
	}

	entry := NewEntry(m.Id, m.Title, m.Status, string(splited[2]))
	post := entry.ConvertFanbox(entry)

	client := fanbox.NewClient(
		&http.Client{},
		config.Default.CreatorId,
		config.Default.SessionId,
		config.Default.CsrfToken,
	)
	client.PushPost(post)

	return nil
}

// YYYY-MM-DD-ID.txt の形で、現在のディレクトリにファイルを保存する
func saveFile(entry Entry) error {
	parsedTime, err := time.Parse(time.RFC3339, entry.updatedAt)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s-%s.md", parsedTime.Format(time.DateOnly), entry.id)

	meta := &meta{
		Id:     entry.id,
		Title:  entry.title,
		Status: string(entry.status),
	}
	metaBytes, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}
	metaString := fmt.Sprintf("---\n%s---\n", string(metaBytes))

	f, err := os.Create(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	_, err = f.Write([]byte(strings.Join([]string{metaString, entry.body}, "\n")))
	if err != nil {
		return err
	}

	return nil
}
