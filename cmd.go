package main

import (
	"log"
	"net/http"

	"github.com/defaultcf/fanboxsync/fanbox"
)

func CommandPull(creator_id string, session_id string) error {
	httpClient := &http.Client{}
	client := fanbox.NewClient(httpClient, creator_id, session_id)
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
