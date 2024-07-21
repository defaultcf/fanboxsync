package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/defaultcf/fanboxsync/fanbox"
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

func CommandCreate(config *config) error {
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

	entry := NewEntry(postId, "", "draft", "")
	entry.updatedAt = time.Now().Format(time.RFC3339)
	err = saveFile(*entry)
	if err != nil {
		return err
	}

	return nil
}

func CommandPush(config *config, path string) error {
	re := regexp.MustCompile(`(\d+)\.md$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) == 0 {
		return fmt.Errorf("can't find id")
	}
	postId := matches[1]
	log.Printf("postId: %s", postId)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	entry := NewEntry(postId, "draft", "draft", string(bytes)) // TODO: タイトルをマークダウンから抽出
	post := entry.ConvertFanbox(entry)
	log.Printf("post: %+v", post)

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

	f, err := os.Create(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	_, err = f.Write([]byte(entry.body))
	if err != nil {
		return err
	}

	return nil
}
