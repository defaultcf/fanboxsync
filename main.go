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
			commandCreate,
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

		err = CommandPull(config)
		return err
	},
}

var commandCreate = &cli.Command{
	Name:  "create",
	Usage: "Create post",
	Action: func(ctx *cli.Context) error {
		log.Print("create")
		config, err := newConfig()
		if err != nil {
			return err
		}

		err = CommandCreate(config)
		return err
	},
}
