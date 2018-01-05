package apiforum

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeShowUsersHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showUsersRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showUsersAction(w, in, db)
	}
	return f
}

type showUsersInput struct {
	Slug string `json:"-"`

	Limit  int64  `json:"limit"`
	Since  string `json:"since"`
	IsDesc bool   `json:"desc"`
}

type showUsersOutputItem api.UserModel

type showUsersOutput []showUsersOutputItem

func showUsersRead(r *http.Request, ps map[string]string) (in showUsersInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func showUsersAction(w http.ResponseWriter, in showUsersInput, db *sql.DB) {
	out := make(showUsersOutput, 0, in.Limit)

	sqlQuery := "SELECT nickname, fullname, email, about FROM users JOIN threads ON (threads.author = users.nickname) WHERE threads.forumSlug = $1 AND nickname > $2 LIMIT $3"
	if in.IsDesc {
		sqlQuery += " ORDER BY nickname DESC"
	} else {
		sqlQuery += " ORDER BY nickname ASC"
	}

	rows, err := db.Query(sqlQuery, in.Slug, in.Since, in.Limit)
	if err != nil {
		log.Println("error: apiforum.showUsersAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showUsersOutputItem
		err = rows.Scan(&outItem.Nickname, &outItem.Fullname, &outItem.Email, &outItem.About)
		if err != nil {
			log.Println("error: apiforum.showUsersAction: SELECT iter:", err)
			w.WriteHeader(500)
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
