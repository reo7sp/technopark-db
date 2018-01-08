package apiservice

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeShowStatusHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		showStatusAction(w, db)
	}
	return f
}

type showStatusOutput api.StatusModel

func showStatusAction(w http.ResponseWriter, db *sql.DB) {
	var out showStatusOutput

	sqlQuery := "SELECT (SELECT count(*) FROM forums), (SELECT count(*) FROM threads), (SELECT count(*) FROM users), (SELECT count(*) FROM posts)"
	err := db.QueryRow(sqlQuery).Scan(&out.ForumsCount, &out.ThreadsCount, &out.UsersCount, &out.PostsCount)
	if err != nil {
		log.Println("error: apiservice.showStatusAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
