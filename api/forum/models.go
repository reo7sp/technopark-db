package forum

import (
	"regexp"
	"database/sql"
)

type ForumModel struct {
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
	PostsCount int `json:"posts"`
	ThreadsCount int `json:"threads"`
}

var slugRegex = regexp.MustCompile("^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$")

func (f *ForumModel) IsValid() bool {
	return slugRegex.MatchString(f.Slug)
}

func (f *ForumModel) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO forum (title, \"user\", slug) VALUES ($1, $2, $3)", f.Title, f.User, f.Slug)
	return err
}

func FindForum(db *sql.DB, slug string) (ForumModel, error) {
	var f ForumModel
	f.Slug = slug
	err := db.QueryRow("SELECT title, \"user\" FROM f WHERE slug=$1", slug).Scan(&f.Title, &f.User)
	if err != nil {
		return f, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE f=$1", f.Slug).Scan(f.PostsCount)
	if err != nil {
		return f, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM threads WHERE f=$1", f.Slug).Scan(f.ThreadsCount)
	if err != nil {
		return f, err
	}

	return f, nil
}
