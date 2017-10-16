package thread

import (
	"database/sql"
)

type ThreadModel struct {
	Id             int32  `json:"id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	ForumSlug      string `json:"forum"`
	Message        string `json:"message"`
	VotesCount     int32  `json:"votes"`
	Slug           string `json:"slug"`
	CreatedDateStr string `json:"created"`
}

func (t *ThreadModel) Create(db *sql.DB) error {
	return db.QueryRow("INSERT INTO thread (title, author, forumSlug, \"message\") VALUES ($1, $2, $3, $4) RETURNING id",
		t.Title, t.Author, t.ForumSlug, t.Message).Scan(&t.Id)
}

func FindThread(db *sql.DB, slug string) (ThreadModel, error) {
	var t ThreadModel
	t.Slug = slug
	err := db.QueryRow("SELECT id, title, author, forumSlug, \"message\", votesCount, createdAt FROM threads WHERE slug=$1",
		slug).Scan(&t.Id, &t.Title, &t.Author, &t.ForumSlug, &t.Message, &t.VotesCount, &t.CreatedDateStr)
	return t, err
}
