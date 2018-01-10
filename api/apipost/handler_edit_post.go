package apipost

import (
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
	"github.com/jackc/pgx"
	"time"
)

func MakeEditPostHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := editPostRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		editPostAction(w, in, db)
	}
	return f
}

type editPostInput struct {
	Id      int64   `json:"-"`
	Message *string `json:"message"`
}

type editPostOutput api.PostModel

func editPostRead(r *http.Request, ps map[string]string) (in editPostInput, err error) {
	id, err := strconv.ParseInt(ps["id"], 10, 64)
	if err != nil {
		return
	}
	in.Id = id

	err = apiutil.ReadJsonObject(r, &in)

	return
}

func editPostAction(w http.ResponseWriter, in editPostInput, db *pgx.ConnPool) {
	var out editPostOutput

	out.Id = in.Id

	sqlQuery := "UPDATE posts SET \"message\" = COALESCE($1, \"message\"), isEdited = ($1 IS NOT NULL AND $1 <> \"message\") WHERE id = $2 RETURNING author::text, createdAt, forumSlug::text, isEdited, threadId, \"message\""
	var t time.Time
	err := db.QueryRow(sqlQuery, in.Message, in.Id).Scan(&out.AuthorNickname, &t, &out.ForumSlug, &out.IsEdited, &out.ThreadId, &out.Message)
	out.CreatedDateStr = t.Format(time.RFC3339Nano)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find post"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apipost.editPostAction: UPDATE:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
