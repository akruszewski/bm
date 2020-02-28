package main

import (
	"log"
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

func main() {
	hp, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	pathDefault := path.Join(hp, ".bm/bm.csv")
	app := cli.App{
		Usage: "manage bookmars from cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: pathDefault,
				Usage: "path to bm database",
			},
		},
		Commands: cli.Commands{
			&cli.Command{
				Name:   "init",
				Usage:  "init bm repository",
				Action: initHandler,
			},
			&cli.Command{
				Name:      "add",
				Usage:     "add bookmark to repository",
				ArgsUsage: "<URL>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "title, t",
						Value: "",
						Usage: "title of the bookmark",
					},
					&cli.StringFlag{
						Name:  "tags",
						Value: "",
						Usage: "tags of the bookmark",
					},
					&cli.StringFlag{
						Name:  "note, n",
						Value: "",
						Usage: "notes to the bookmark",
					},
				},
				Action: addHandler,
			},
			&cli.Command{
				Name:      "update",
				Usage:     "update bookmark",
				ArgsUsage: "<URL>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "title, t",
						Value: "NOT_UPDATED",
						Usage: "title of the bookmark",
					},
					&cli.StringFlag{
						Name:  "tags",
						Value: "NOT_UPDATED",
						Usage: "tags of the bookmark",
					},
					&cli.StringFlag{
						Name:  "note, n",
						Value: "NOT_UPDATED",
						Usage: "notes to the bookmark",
					},
				},
				Action: updateHandler,
			},
			&cli.Command{
				Name:  "list",
				Usage: "list bookmarks",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "fields",
						Value: "url",
						Usage: "fields which will be displayed",
					},
					&cli.StringFlag{
						Name:  "tags, t",
						Value: "",
						Usage: "tags of displayed bookmarks",
					},
				},
				Action: listHandler,
			},
			&cli.Command{
				Name:      "get",
				Usage:     "get bookmark",
				ArgsUsage: "<URL>",
				Action:    getHandler,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
