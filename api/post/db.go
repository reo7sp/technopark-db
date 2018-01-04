package post

import "database/sql"

func CreatePostInDB(p *PostModel, db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO post (title, author, forumSlug, \"message\") VALUES ($1, $2, $3, $4) RETURNING id",
		p.Title, p.AuthorNickname, p.ForumSlug, p.Message,
	).Scan(&p.Id)
	return err
}

func EditMessageOfPostByIdInDB(id int64, newMessage string, db *sql.DB) (int, error) {

}

func LoadUpPostHavingIdAndMessageFromDB(p *PostModel, db *sql.DB) error {

}

func FindPostByIdInDB(id int64, db *sql.DB) (PostModel, error) {
	var t PostModel
	t.Id = id
	err := db.QueryRow(
		"SELECT slug, title, author, forumSlug, \"message\", votesCount, createdAt FROM posts WHERE id=$1",
		id,
	).Scan(&t.Slug, &t.Title, &t.AuthorNickname, &t.ForumSlug, &t.Message, &t.VotesCount, &t.CreatedDateStr)
	return t, err
}
