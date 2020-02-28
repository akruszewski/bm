package main

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const comma = '|'

var hd = []string{
	"title",
	"url",
	"tags",
	"notes",
	"document",
	"created_at",
	"updated_at",
}

type Bookmark struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
	Notes string   `json:"notes"`

	Document  string    `json:"document"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func initRepo(repoPath string) error {
	if err := ioutil.WriteFile(
		repoPath,
		[]byte(strings.Join(hd, string(comma))+"\n"),
		0644,
	); err != nil {
		return err
	}
	return nil
}

func add(repoPath string, bm *Bookmark) error {
	f, err := os.OpenFile(
		repoPath,
		os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close()
	w := newWriter(f)
	if err = w.Write(marshal(bm)); err != nil {
		return err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

//TODO: this implementation is shitty... rewrite this function
func update(repoPath string, bm *Bookmark) error {
	f, err := os.OpenFile(
		repoPath,
		os.O_RDWR,
		0644,
	)

	if err != nil {
		return err
	}
	defer f.Close()
	r := newReader(f)
	w := newWriter(f)

	lines, err := r.ReadAll()
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	//Add header
	if err = w.Write(lines[0]); err != nil {
		return err
	}
	for _, line := range lines[1:] {
		b, err := parseBookmark(line)
		if err != nil {
			return err
		}
		if bm.URL == b.URL {
			b = updateBookmark(b, bm)

			if err = w.Write(marshal(b)); err != nil {
				return err
			}
		} else {
			if err = w.Write(line); err != nil {
				return err
			}
		}

	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func get(repoPath, url string) (*Bookmark, error) {
	f, err := os.Open(repoPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := newReader(f)

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	if !validate(header) {
		return nil, errors.New("invalid csv file")
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		bm, err := parseBookmark(record)
		if err != nil {
			return nil, err
		}
		if bm.URL == url {
			return bm, nil
		}
	}
	return nil, errors.New("bookmark not found")
}

func list(repoPath string, tags ...string) ([]*Bookmark, error) {
	f, err := os.Open(repoPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tagMap := make(map[string]bool)
	for _, t := range tags {
		if t != "" {
			tagMap[t] = true
		}
	}
	tmLen := len(tagMap)

	r := newReader(f)
	bms := []*Bookmark{}

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	if !validate(header) {
		return nil, errors.New("invalid csv file")
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		bm, err := parseBookmark(record)
		if err != nil {
			return nil, err
		}
		if tmLen == 0 || anyIn(tagMap, bm.Tags) {
			bms = append(bms, bm)
		}
	}
	return bms, nil
}

func anyIn(td map[string]bool, ts2 []string) bool {
	for _, t := range ts2 {
		if _, ok := td[t]; ok {
			return true
		}
	}
	return false
}

func validate(header []string) bool {
	for i, h := range header {
		if h != hd[i] {
			return false
		}
	}
	return true
}

func updateBookmark(bm, bmPatch *Bookmark) *Bookmark {
	if bmPatch.Title != "NOT_UPDATED" {
		bm.Title = bmPatch.Title
	}
	if bmPatch.Tags[0] != "NOT_UPDATED" {
		bm.Tags = bmPatch.Tags
	}
	if bmPatch.Notes != "NOT_UPDATED" {
		bm.Notes = bmPatch.Notes
	}
	bm.UpdatedAt = time.Now().UTC()
	return bm
}

func marshal(bm *Bookmark) []string {
	return []string{
		bm.Title,
		bm.URL,
		strings.Join(bm.Tags, ";"),
		bm.Notes,
		bm.Document,
		bm.CreatedAt.Format(time.RFC3339),
		bm.UpdatedAt.Format(time.RFC3339),
	}
}

func parseBookmark(data []string) (*Bookmark, error) {
	tags := strings.Split(data[2], ";")
	cr, err := time.Parse(time.RFC3339, data[5])
	if err != nil {
		return nil, err
	}
	up, err := time.Parse(time.RFC3339, data[6])
	if err != nil {
		return nil, err
	}
	return &Bookmark{
		Title:     data[0],
		URL:       data[1],
		Tags:      tags,
		Notes:     data[3],
		Document:  data[4],
		CreatedAt: cr,
		UpdatedAt: up,
	}, nil
}

func newWriter(r io.Writer) *csv.Writer {
	w := csv.NewWriter(r)
	w.Comma = comma
	return w
}

func newReader(r io.Reader) *csv.Reader {
	w := csv.NewReader(r)
	w.Comma = comma
	return w
}
