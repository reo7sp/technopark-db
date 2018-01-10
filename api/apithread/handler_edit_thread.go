package apithread

import (
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"github.com/jackc/pgx"
	"time"
)

func MakeEditThreadHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := editThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		editThreadAction(w, in, db)
	}
	return f
}

type editThreadInput struct {
	slugOrIdInput

	Title   *string `json:"title"`
	Message *string `json:"message"`
}

type editThreadOutput api.ThreadModel

func editThreadRead(r *http.Request, ps map[string]string) (in editThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func editThreadAction(w http.ResponseWriter, in editThreadInput, db *pgx.ConnPool) {
	var out editThreadOutput

	sqlQuery := "UPDATE threads SET title = COALESCE($1, title), \"message\" = COALESCE($2, \"message\")"
	sqlFields := []interface{}{in.Title, in.Message, nil}
	if in.HasId {
		sqlQuery += " WHERE id = $3"
		sqlFields[2] = in.Id
	} else {
		sqlQuery += " WHERE slug = $3"
		sqlFields[2] = in.Slug
	}
	sqlQuery += " RETURNING author, createdAt, forumSlug, id, \"message\", slug, title"

	var t time.Time
	err := db.QueryRow(sqlQuery, sqlFields...).Scan(&out.AuthorNickname, &t, &out.ForumSlug, &out.Id, &out.Message, &out.Slug, &out.Title)
	out.CreatedDateStr = t.Format(time.RFC3339Nano)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.editThreadAction: UPDATE:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
