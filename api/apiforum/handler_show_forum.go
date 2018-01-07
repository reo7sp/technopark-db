package apiforum

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"errors"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/dbutil"
)

func MakeShowForumHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showForumRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showForumAction(w, in, db)
	}
	return f
}

type showForumInput struct {
	Slug string
}

type showForumOutput api.ForumModel

func showForumRead(r *http.Request, ps map[string]string) (in showForumInput, err error) {
	slug, ok := ps["slug"]
	in.Slug = slug
	if !ok {
		err = errors.New("slug is empty")
	}
	return
}

func showForumAction(w http.ResponseWriter, in showForumInput, db *sql.DB) {
	var out showForumOutput

	sqlQuery := "SELECT slug, title, \"user\", postsCount, threadsCount FROM forums WHERE slug = $1"
	err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Slug, &out.Title, &out.User, &out.PostsCount, &out.ThreadsCount)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.Error{Message: "Can't find forum"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apiforum.showForumAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
