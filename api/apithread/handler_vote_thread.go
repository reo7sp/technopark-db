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

func MakeVoteThreadHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := voteThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		voteThreadAction(w, in, db)
	}
	return f
}

type voteThreadInput struct {
	slugOrIdInput

	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
}

type voteThreadOutput api.ThreadModel

func voteThreadRead(r *http.Request, ps map[string]string) (in voteThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func voteThreadAction(w http.ResponseWriter, in voteThreadInput, db *pgx.ConnPool) {
	var out voteThreadOutput

	sqlQuery := `
	INSERT INTO votes (nickname, threadId, voice) VALUES (
		$4,
		(
			CASE WHEN $1 IS TRUE
			THEN $2
			ELSE (SELECT id FROM threads WHERE slug = $3::citext)
			END
		),
		$5
	)
	ON CONFLICT (nickname, threadId) DO UPDATE SET voice = EXCLUDED.voice;
    `

	_, err := db.Exec(sqlQuery, in.HasId, in.Id, in.Slug, in.Nickname, in.Voice)

	if err != nil && dbutil.IsErrorAboutFailedForeignKey(err) {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.voteThreadAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	sqlQuery = `
	SELECT id, title, author::text, forumSlug::text, "message", votesCount, slug::text, createdAt FROM threads
	WHERE (
		CASE WHEN $1 IS TRUE
		THEN id = $2
		ELSE slug = $3::citext
		END
	)
	`
	var t time.Time
	err = db.QueryRow(sqlQuery, in.HasId, in.Id, in.Slug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.ForumSlug, &out.Message, &out.VotesCount, &out.Slug, &t)
	out.CreatedDateStr = t.Format(time.RFC3339Nano)
	if err != nil {
		log.Println("error: apithread.voteThreadAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
