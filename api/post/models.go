package post

import (
	"github.com/reo7sp/technopark-db/api/thread"
	"github.com/reo7sp/technopark-db/api/forum"
	"github.com/reo7sp/technopark-db/api/user"
)

type PostModel struct {
	Id              int64  `json:"id"`
	ParentMessageId int64  `json:"parent"`
	AuthorNickname  string `json:"author"`
	Message         string `json:"message"`
	IsEdited        bool   `json:"isEdited"`
	ForumSlug       string `json:"forum"`
	ThreadId        int64  `json:"thread"`
	CreatedDateStr  string `json:"created"`
}

type PostFullModel struct {
	Post   PostModel          `json:"post"`
	Author user.UserModel     `json:"author"`
	Thread thread.ThreadModel `json:"thread"`
	Forum  forum.ForumModel   `json:"forum"`
}
