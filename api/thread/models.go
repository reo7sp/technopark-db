package thread

type ThreadModel struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	ForumSlug      string `json:"forum"`
	Message        string `json:"message"`
	VotesCount     int64  `json:"votes"`
	Slug           string `json:"slug"`
	CreatedDateStr string `json:"created"`
}

type VoteModel struct {
	Nickname string `json:"nickname"`
	Voice    int64   `json:"voice"`
}
