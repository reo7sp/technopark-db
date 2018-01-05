package apiforum

import (
	"database/sql"
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeCreateForumHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createForumRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createForumAction(w, in, db)
	}
	return f
}

type createForumInput struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type createForumOutput api.ForumModel

func createForumRead(r *http.Request, ps map[string]string) (in createForumInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func createForumAction(w http.ResponseWriter, in createForumInput, db *sql.DB) {
	var out createForumOutput

	out.Slug = in.Slug

	sqlQuery := "INSERT INTO forums (title, \"user\", slug) VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlQuery, in.Title, in.User, in.Slug)

	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT title, \"user\", postsCount, threadsCount FROM forums WHERE slug = $1"
		err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Title, &out.User, &out.PostsCount, &out.ThreadsCount)

		if err != nil {
			log.Println("error: apiforum.createForumAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, out, 409)
		return
	}
	if err != nil {
		log.Println("error: apiforum.createForumAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	out.Title = in.Title
	out.User = in.User
	out.PostsCount = 0
	out.ThreadsCount = 0

	apiutil.WriteJsonObject(w, out, 200)
}
