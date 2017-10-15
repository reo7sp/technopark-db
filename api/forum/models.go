package forum

import (
	"regexp"
	"database/sql"
)

type Forum struct {
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
	PostsCount int `json:"posts"`
	ThreadsCount int `json:"threads"`
}

var slugRegex = regexp.MustCompile("^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$")

func (f *Forum) IsValid() bool {
	return slugRegex.MatchString(f.Slug)
}

func (f *Forum) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO forum (title, \"user\", slug) VALUES ($1, $2, $3)", f.Title, f.User, f.Slug)
	return err
}

func (f *Forum) FetchPostsAndThreadsCount(db *sql.DB) error {
	panic("Not implemented")
}
