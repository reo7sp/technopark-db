package thread

import "database/sql"

func CreateThreadInDB(t ThreadModel, db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO thread (title, author, forumSlug, \"message\") VALUES ($1, $2, $3, $4) RETURNING id",
		t.Title, t.Author, t.ForumSlug, t.Message,
	).Scan(&t.Id)
	return err
}

func FindThreadByIdInDB(id int64, db *sql.DB) (ThreadModel, error) {
	var t ThreadModel
	t.Id = id
	err := db.QueryRow(
		"SELECT slug, title, author, forumSlug, \"message\", votesCount, createdAt FROM threads WHERE id=$1",
		id,
	).Scan(&t.Slug, &t.Title, &t.Author, &t.ForumSlug, &t.Message, &t.VotesCount, &t.CreatedDateStr)
	return t, err
}

func FindThreadBySlugInDB(slug string, db *sql.DB) (ThreadModel, error) {
	var t ThreadModel
	t.Slug = slug
	err := db.QueryRow(
		"SELECT id, title, author, forumSlug, \"message\", votesCount, createdAt FROM threads WHERE slug=$1",
		slug,
	).Scan(&t.Id, &t.Title, &t.Author, &t.ForumSlug, &t.Message, &t.VotesCount, &t.CreatedDateStr)
	return t, err
}

func EditVotesCountOfThreadInDB(t ThreadModel, db *sql.DB) error {

}