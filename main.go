package main

import (
	"log"
	"os"

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
		log.Printf("creator_id: %s", config.Default.CreatorId)
		log.Printf("session_id: %s", config.Default.SessionId)

		return nil
	},
}
