package apipost

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"strconv"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/dbutil"
)

func MakeEditPostHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
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

func editPostAction(w http.ResponseWriter, in editPostInput, db *sql.DB) {
	var out editPostOutput

	out.Id = in.Id

	sqlQuery := "UPDATE posts SET \"message\" = COALESCE($1, \"message\"), isEdited = ($1 IS NOT NULL AND $1 <> \"message\") WHERE id = $2 RETURNING author, createdAt, forumSlug, isEdited, threadId, \"message\""
	err := db.QueryRow(sqlQuery, in.Message, in.Id).Scan(&out.AuthorNickname, &out.CreatedDateStr, &out.ForumSlug, &out.IsEdited, &out.ThreadId, &out.Message)

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
