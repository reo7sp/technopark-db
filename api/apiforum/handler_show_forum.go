package apiforum

import (
	"errors"
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
)

func MakeShowForumHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
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

func showForumAction(w http.ResponseWriter, in showForumInput, db *pgx.ConnPool) {
	var out showForumOutput

	sqlQuery := "SELECT slug::text, title, \"user\"::text, postsCount, threadsCount FROM forums WHERE slug = $1::citext"
	err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Slug, &out.Title, &out.User, &out.PostsCount, &out.ThreadsCount)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find forum"}
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
