package apiuser

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeShowUserHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showUserRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showUserAction(w, in, db)
	}
	return f
}

type showUserInput struct {
	Nickname string
}

type showUserOutput api.UserModel

func showUserRead(r *http.Request, ps map[string]string) (in showUserInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func showUserAction(w http.ResponseWriter, in showUserInput, db *sql.DB) {
	var out showUserOutput

	out.Nickname = in.Nickname

	sqlQuery := "SELECT fullname, about, email FROM users WHERE nickname = $1"
	err := db.QueryRow(sqlQuery, in.Nickname).Scan(&out.Fullname, &out.About, &out.Email)

	if err != nil {
		log.Println("error: apiuser.showUserAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
