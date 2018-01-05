package apipost

type PostModel struct {
	Id             int64  `json:"id"`
	ParentPostId   int64  `json:"parent"`
	AuthorNickname string `json:"author"`
	Message        string `json:"message"`
	IsEdited       bool   `json:"isEdited"`
	ForumSlug      string `json:"forum"`
	ThreadId       int64  `json:"thread"`
	CreatedDateStr string `json:"created"`
}
