package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/defaultcf/fanboxsync/fanbox"
)

func CommandPull(creatorId string, sessionId string) error {
	client := fanbox.NewClient(&http.Client{}, creatorId, sessionId)
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

// YYYY-MM-DD-ID.txt の形で、現在のディレクトリにファイルを保存する
func saveFile(entry Entry) error {
	parsedTime, err := time.Parse(time.RFC3339, entry.updatedAt)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s-%s.txt", parsedTime.Format(time.DateOnly), entry.id)

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
