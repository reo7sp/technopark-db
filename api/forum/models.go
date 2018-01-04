package forum

import (
	"regexp"
)

type ForumModel struct {
	Title        string `json:"title"`
	User         string `json:"user"`
	Slug         string `json:"slug"`
	PostsCount   int    `json:"posts"`
	ThreadsCount int    `json:"threads"`
}

var slugRegex = regexp.MustCompile("^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$")

func (f *ForumModel) IsValid() bool {
	return slugRegex.MatchString(f.Slug)
}
