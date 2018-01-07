package apiforum

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api"
	"strconv"
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

type showThreadsOutputItem api.ThreadModel

type showThreadsOutput []showThreadsOutputItem

func showThreadsRead(r *http.Request, ps map[string]string) (in showThreadsInput, err error) {
	in.Slug = ps["slug"]

	query := r.URL.Query()

	in.Limit, err = strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		err = nil
		in.Limit = -1
	}
	in.Since = query.Get("since")
	in.IsDesc = query.Get("desc") == "true"

	return
}

func showThreadsAction(w http.ResponseWriter, in showThreadsInput, db *sql.DB) {
	var out showThreadsOutput

	if in.Limit != -1 {
		out = make(showThreadsOutput, 0, in.Limit)
	} else {
		out = make(showThreadsOutput, 0)
	}

	forumSlug := ""
	err := db.QueryRow("SELECT slug FROM forums WHERE slug = $1", in.Slug).Scan(&forumSlug)
	if err != nil {
		errJson := api.Error{Message: "Can't find forum"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	sqlFields := make([]interface{}, 0, 3)
	sqlFields = append(sqlFields, in.Slug)
	sqlQuery := "SELECT id, title, author, \"message\", createdAt, votesCount, slug, forumSlug FROM threads WHERE forumSlug = $1"
	if in.IsDesc {
		if in.Since != "" {
			sqlFields = append(sqlFields, in.Since)
			sqlQuery += " AND createdAt <= $2"
		}
		sqlQuery += " ORDER BY createdAt DESC"
	} else {
		if in.Since != "" {
			sqlFields = append(sqlFields, in.Since)
			sqlQuery += " AND createdAt >= $2"
		}
		sqlQuery += " ORDER BY createdAt ASC"
	}
	if in.Limit != -1 {
		sqlFields = append(sqlFields, in.Limit)
		sqlQuery += " LIMIT $" + strconv.FormatInt(int64(len(sqlFields)), 10)
	}

	rows, err := db.Query(sqlQuery, sqlFields...)
	if err != nil {
		log.Println("error: apiforum.showThreadsAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showThreadsOutputItem
		err = rows.Scan(&outItem.Id, &outItem.Title, &outItem.AuthorNickname, &outItem.Message, &outItem.CreatedDateStr, &outItem.VotesCount, &outItem.Slug, &outItem.ForumSlug)
		if err != nil {
			log.Println("error: apiforum.showThreadsAction: SELECT iter:", err)
			w.WriteHeader(500)
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
