package apithread

import (
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"time"
)

func MakeShowThreadHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showThreadAction(w, in, db)
	}
	return f
}

type showThreadInput struct {
	slugOrIdInput
}

type showThreadOutput api.ThreadModel

func showThreadRead(r *http.Request, ps map[string]string) (in showThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	return
}

func showThreadAction(w http.ResponseWriter, in showThreadInput, db *pgx.ConnPool) {
	var out showThreadOutput

	sqlQuery := `
	SELECT id, title, author::text, forumSlug::text, "message", votesCount, slug::text, createdAt FROM threads
	WHERE (
		CASE WHEN $1 IS TRUE
		THEN id = $2
		ELSE slug = $3::citext
		END
	)
	`
	var t time.Time
	err := db.QueryRow(sqlQuery, in.HasId, in.Id, in.Slug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.ForumSlug, &out.Message, &out.VotesCount, &out.Slug, &t)
	out.CreatedDateStr = t.UTC().Format(api.TIMEFORMAT)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.showThreadAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
