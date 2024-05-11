package main

import (
	"log"
	"os"

	"github.com/defaultcf/fanboxsync/fanbox"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "fanboxsync",
		Usage:   "Sync FANBOX posts",
		Version: version,
		Commands: []*cli.Command{
			commandPull,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var commandPull = &cli.Command{
	Name:  "pull",
	Usage: "Pull posts from FANBOX",
	Action: func(ctx *cli.Context) error {
		log.Print("pull")
		config, err := newConfig()
		if err != nil {
			return err
		}
		creator_id := config.Default.CreatorId
		session_id := config.Default.SessionId

		client := fanbox.NewClient()
		posts, err := client.GetPosts(creator_id, session_id)
		if err != nil {
			return err
		}
		for _, v := range posts {
			log.Printf("post: %+v\n", v)
		}

		return nil
	},
}
