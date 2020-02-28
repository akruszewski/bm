package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func initHandler(c *cli.Context) error {
	return initRepo(c.String("path"))

}

func getHandler(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("URL argument required")
	}
	bm, err := get(c.String("path"), c.Args().First())
	fmt.Printf(`Title:     %s
URL:       %s
Tags:      %s
CreateAt:  %s
UpdatedAt: %s
Notes:
	%s
`, bm.Title, bm.URL, bm.Tags, bm.CreatedAt, bm.UpdatedAt, bm.Notes)

	if err != nil {
		return err
	}
	return nil
}

func addHandler(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("URL argument required")
	}
	return add(c.String("path"), &Bookmark{
		Title:     c.String("title"),
		URL:       c.Args().First(),
		Tags:      strings.Split(c.String("tags"), ";"),
		Notes:     c.String("note"),
		Document:  ".",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}
func updateHandler(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("URL argument required")
	}
	return update(c.String("path"), &Bookmark{
		Title:    c.String("title"),
		URL:      c.Args().First(),
		Tags:     strings.Split(c.String("tags"), ";"),
		Notes:    c.String("note"),
		Document: ".",
	})
}

func listHandler(c *cli.Context) error {
	var err error
	var bms []*Bookmark
	tgs := strings.Split(c.String("tags"), ";")
	repoPath := c.String("path")

	if len(tgs) == 0 {
		bms, err = list(repoPath)
	} else {
		bms, err = list(repoPath, tgs...)

	}
	if err != nil {
		return err
	}
	fields := c.String("fields")
	for _, bm := range bms {
		if strings.Contains(fields, "title") {
			fmt.Printf("%s\t", bm.Title)
		}
		if strings.Contains(fields, "url") {
			fmt.Printf("%s\t", bm.URL)
		}
		if strings.Contains(fields, "tags") {
			fmt.Printf("%s", strings.Join(bm.Tags, ","))
		}
		fmt.Printf("\n")
	}
	return nil
}
