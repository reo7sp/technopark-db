package apiforum

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"errors"
	"github.com/reo7sp/technopark-db/dbutil"
	"github.com/reo7sp/technopark-db/api"
)

func MakeCreateThreadHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createThreadAction(w, in, db)
	}
	return f
}

type createThreadInput struct {
	ForumSlug string `json:"-"`

	Title   string `json:"title"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type createThreadOutput api.ThreadModel

func createThreadRead(r *http.Request, ps map[string]string) (in createThreadInput, err error) {
	slug, ok := ps["slug"]
	in.ForumSlug = slug
	if !ok {
		err = errors.New("slug is empty")
		return
	}

	err = apiutil.ReadJsonObject(r, &in)

	return
}

func createThreadAction(w http.ResponseWriter, in createThreadInput, db *sql.DB) {
	var out createThreadOutput

	out.Slug = in.ForumSlug

	sqlQuery := "INSERT INTO threads (title, author, forumSlug, \"message\") VALUES ($1, $2, $3, $4) RETURNING id, createdAt"
	err := db.QueryRow(sqlQuery, in.Title, in.Author, in.ForumSlug, in.Message).Scan(&out.Id, &out.CreatedDateStr)

	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT id, title, author, \"message\", votes, createdAt FROM threads WHERE slug = $1"
		err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.Message, &out.VotesCount, &out.CreatedDateStr)

		if err != nil {
			log.Println("error: apiforum.createThreadAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, out, 409)
		return
	}
	if err != nil {
		log.Println("error: apiforum.createThreadAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	out.Title = in.Title
	out.AuthorNickname = in.Author
	out.Message = in.Message
	out.VotesCount = 0

	apiutil.WriteJsonObject(w, out, 200)
}
