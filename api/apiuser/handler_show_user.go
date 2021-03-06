package apiuser

import (
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
)

func MakeShowUserHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
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
	in.Nickname = ps["nickname"]
	return
}

func showUserAction(w http.ResponseWriter, in showUserInput, db *pgx.ConnPool) {
	var out showUserOutput

	sqlQuery := "SELECT nickname::text, fullname, about, email::text FROM users WHERE nickname = $1::citext"
	err := db.QueryRow(sqlQuery, in.Nickname).Scan(&out.Nickname, &out.Fullname, &out.About, &out.Email)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find user"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apiuser.showUserAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
