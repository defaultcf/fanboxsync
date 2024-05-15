package main

import (
	"log"
	"net/http"

	"github.com/defaultcf/fanboxsync/fanbox"
)

func CommandPull(creatorId string, sessionId string) error {
	httpClient := &http.Client{}
	client := fanbox.NewClient(httpClient, creatorId, sessionId)
	posts, err := client.GetPosts()
	if err != nil {
		return err
	}
	for _, v := range posts {
		post, err := client.GetPost(v.Id)
		if err != nil {
			return err
		}
		//log.Printf("post: %+v\n", post)

		e := NewEntry("", "", "", "")
		e.ConvertPost(post)
		log.Printf("entry: %+v\n", e)
	}
	return nil
}
