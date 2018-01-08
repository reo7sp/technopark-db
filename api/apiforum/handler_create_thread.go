package apiforum

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"errors"
	"github.com/reo7sp/technopark-db/api"
	"time"
	"github.com/reo7sp/technopark-db/dbutil"
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

	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Message      string  `json:"message"`
	ThreadSlug   *string `json:"slug"`
	CreatedAtStr string  `json:"created"`
}

type createThreadOutput api.ThreadModel

type createThreadGetForumInfo struct {
	ForumSlug string
}

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

func createThreadGetForum(in createThreadInput, db *sql.DB) (r createThreadGetForumInfo, err error) {
	sqlQuery := "SELECT slug FROM forums WHERE slug = $1"
	err = db.QueryRow(sqlQuery, in.ForumSlug).Scan(&r.ForumSlug)
	return
}

func createThreadAction(w http.ResponseWriter, in createThreadInput, db *sql.DB) {
	forumInfo, err := createThreadGetForum(in, db)
	if err != nil {
		errJson := api.ErrorModel{Message: "Can't find forum"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	var out createThreadOutput

	out.ForumSlug = forumInfo.ForumSlug

	if in.CreatedAtStr == "" {
		in.CreatedAtStr = time.Now().Format(time.RFC3339)
	}

	sqlQuery := "INSERT INTO threads (slug, title, author, forumSlug, \"message\", createdAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = db.QueryRow(sqlQuery, in.ThreadSlug, in.Title, in.Author, in.ForumSlug, in.Message, in.CreatedAtStr).Scan(&out.Id)

	if err != nil && dbutil.IsErrorAboutFailedForeignKey(err) {
		errJson := api.ErrorModel{Message: "Can't find user"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT id, title, author, \"message\", createdAt, slug, forumSlug, createdAt FROM threads WHERE slug = $1"
		err := db.QueryRow(sqlQuery, in.ThreadSlug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.Message, &out.CreatedDateStr, &out.Slug, &out.ForumSlug, &out.CreatedDateStr)

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
	out.CreatedDateStr = in.CreatedAtStr
	if in.ThreadSlug != nil {
		out.Slug = *in.ThreadSlug
	}

	apiutil.WriteJsonObject(w, out, 201)
}
