package api

import (
	"regexp"
)

type ErrorModel struct {
	Message string
}

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

type PostModel struct {
	Id             int64  `json:"id"`
	ParentPostId   *int64 `json:"parent"`
	AuthorNickname string `json:"author"`
	Message        string `json:"message"`
	IsEdited       bool   `json:"isEdited"`
	ForumSlug      string `json:"forum"`
	ThreadId       int64  `json:"thread"`
	CreatedDateStr string `json:"created"`
}

type StatusModel struct {
	UsersCount   int64 `json:"user"`
	ForumsCount  int64 `json:"forum"`
	ThreadsCount int64 `json:"thread"`
	PostsCount   int64 `json:"post"`
}

type ThreadModel struct {
	Id             int64   `json:"id"`
	Title          string  `json:"title"`
	AuthorNickname string  `json:"author"`
	ForumSlug      string  `json:"forum"`
	Message        string  `json:"message"`
	VotesCount     int64   `json:"votes"`
	Slug           *string `json:"slug"`
	CreatedDateStr string  `json:"created"`
}

type UserModel struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
