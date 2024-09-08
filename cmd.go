package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/goccy/go-yaml"
)

func userAgent() string {
	return fmt.Sprintf("fanboxsync/%s", version)
}

func CommandPull(config *config) error {
	f, err := fanbox.NewFanbox(config.Default.CsrfToken, config.Default.SessionId, userAgent())
	if err != nil {
		return err
	}

	posts, err := f.GetPosts()
	if err != nil {
		return err
	}

	for _, v := range posts {
		post, err := f.GetPost(v.ID.Value)
		if err != nil {
			return err
		}

		e := NewEntry("", "", "", "", "")
		converted := e.ConvertPost(&post)
		fmt.Printf("%+v\n", converted)

		err = saveFile(*converted)
		if err != nil {
			return err
		}
	}

	return nil
}

func CommandCreate(config *config, title string) error {
	f, err := fanbox.NewFanbox(config.Default.CsrfToken, config.Default.SessionId, userAgent())
	if err != nil {
		return err
	}

	postId, err := f.CreatePost()
	if err != nil {
		return err
	}

	entry := NewEntry(postId, title, "draft", "0", "")
	post := entry.ConvertFanbox(entry)
	_, err = f.PushPost(post) // タイトルをセット
	if err != nil {
		return err
	}

	entry.UpdatedAt = time.Now().Format(time.RFC3339)
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
	Fee    string `yaml:"fee"`
}

func CommandPush(config *config, path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	bytes, err := io.ReadAll(fi)
	if err != nil {
		return err
	}

	// メタデータをマークダウンから抽出
	rawBody := string(bytes)
	reMeta := regexp.MustCompile(`---\n\n`)
	splited := reMeta.Split(rawBody, 2)
	m := meta{}
	err = yaml.Unmarshal([]byte(splited[0]), &m)
	if err != nil {
		return err
	}

	entry := NewEntry(m.Id, m.Title, m.Status, m.Fee, string(splited[1]))
	post := entry.ConvertFanbox(entry)

	f, err := fanbox.NewFanbox(config.Default.CsrfToken, config.Default.SessionId, userAgent())
	if err != nil {
		return err
	}

	_, err = f.PushPost(post)
	if err != nil {
		return err
	}

	return nil
}

func CommandDelete(config *config, path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	bytes, err := io.ReadAll(fi)
	if err != nil {
		return err
	}

	rawBody := string(bytes)
	reMeta := regexp.MustCompile(`---\n\n`)
	splited := reMeta.Split(rawBody, 2)
	m := meta{}
	err = yaml.Unmarshal([]byte(splited[0]), &m)
	if err != nil {
		return err
	}

	f, err := fanbox.NewFanbox(config.Default.CsrfToken, config.Default.SessionId, userAgent())
	if err != nil {
		return err
	}
	err = f.DeletePost(m.Id)
	if err != nil {
		return err
	}

	return nil
}

// YYYY-MM-DD-ID.txt の形で、現在のディレクトリにファイルを保存する
func saveFile(entry Entry) error {
	parsedTime, err := time.Parse(time.RFC3339, entry.UpdatedAt)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s-%s.md", parsedTime.Format(time.DateOnly), entry.ID)

	meta := &meta{
		Id:     entry.ID,
		Title:  entry.Title,
		Status: string(entry.Status),
		Fee:    string(entry.Fee),
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

	_, err = f.Write([]byte(strings.Join([]string{metaString, entry.Body}, "\n")))
	if err != nil {
		return err
	}

	return nil
}
