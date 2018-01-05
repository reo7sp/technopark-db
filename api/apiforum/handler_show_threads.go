package apiforum

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api/apithread"
)

func MakeShowThreadsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showThreadsRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showThreadsAction(w, in, db)
	}
	return f
}

type showThreadsInput struct {
	Slug string `json:"-"`

	Limit  int64  `json:"limit"`
	Since  string `json:"since"`
	IsDesc bool   `json:"desc"`
}

type showThreadsOutputItem apithread.ThreadModel

type showThreadsOutput []showThreadsOutputItem

func showThreadsRead(r *http.Request, ps map[string]string) (in showThreadsInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func showThreadsAction(w http.ResponseWriter, in showThreadsInput, db *sql.DB) {
	out := make(showThreadsOutput, 0, in.Limit)

	sqlQuery := "SELECT id, title, author, \"message\", votes, createdAt FROM threads WHERE forumSlug = $1 AND createdAt >= $2 LIMIT $3"
	if in.IsDesc {
		sqlQuery += " ORDER BY createdAt DESC"
	} else {
		sqlQuery += " ORDER BY createdAt ASC"
	}

	rows, err := db.Query(sqlQuery, in.Slug, in.Since, in.Limit)
	if err != nil {
		log.Println("error: apiforum.showThreadsAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showThreadsOutputItem
		outItem.ForumSlug = in.Slug
		err = rows.Scan(&outItem.Id, &outItem.Title, &outItem.AuthorNickname, &outItem.Message, &outItem.VotesCount, &outItem.CreatedDateStr)
		if err != nil {
			log.Println("error: apiforum.showThreadsAction: SELECT iter:", err)
			w.WriteHeader(500)
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
