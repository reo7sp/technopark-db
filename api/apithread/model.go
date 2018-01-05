package apithread

type ThreadModel struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	AuthorNickname string `json:"author"`
	ForumSlug      string `json:"forum"`
	Message        string `json:"message"`
	VotesCount     int64  `json:"votes"`
	Slug           string `json:"slug"`
	CreatedDateStr string `json:"created"`
}
