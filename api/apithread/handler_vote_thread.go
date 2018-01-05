package apithread

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeVoteThreadHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
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

type voteThreadGetThreadInfo struct {
	Id int64
}

func voteThreadRead(r *http.Request, ps map[string]string) (in voteThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func voteThreadGetThread(in voteThreadInput, db *sql.DB) (r voteThreadGetThreadInfo, err error) {
	if !in.HasId {
		sqlQuery := "SELECT id FROM threads WHERE slug = $1"
		err = db.QueryRow(sqlQuery, in.Slug).Scan(&r.Id)
	}
	return
}

func voteThreadAction(w http.ResponseWriter, in voteThreadInput, db *sql.DB) {
	r, err := voteThreadGetThread(in, db)
	if err != nil {
		log.Println("error: apithread.voteThreadAction: voteThreadGetThread:", err)
		w.WriteHeader(500)
		return
	}

	var out voteThreadOutput

	sqlQuery := `
	INSERT INTO votes (nickname, threadId, voice) VALUES ($1, $2, $3)
	ON CONFLICT (nickname, threadId) DO UPDATE SET voice = excluded.voice
    `

	_, err = db.Exec(sqlQuery, in.Nickname, r.Id, in.Voice)
	if err != nil {
		log.Println("error: apithread.voteThreadAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
