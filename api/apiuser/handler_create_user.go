package apiuser

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/dbutil"
	"github.com/reo7sp/technopark-db/api"
	"github.com/lib/pq"
)

func MakeCreateUserHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createUserRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createUserAction(w, in, db)
	}
	return f
}

type createUserInput struct {
	Nickname string `json:"-"`

	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type createUserOutput api.UserModel

func createUserRead(r *http.Request, ps map[string]string) (in createUserInput, err error) {
	in.Nickname = ps["nickname"]
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func createUserAction(w http.ResponseWriter, in createUserInput, db *sql.DB) {
	var out createUserOutput

	sqlQuery := "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlQuery, in.Nickname, in.Fullname, in.About, in.Email)

	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT fullname, about, email FROM users WHERE nickname = $1"
		err := db.QueryRow(sqlQuery, in.Nickname).Scan(&out.Fullname, &out.About, &out.Email)

		if err != nil {
			log.Println("error: apiforum.createForumAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, out, 409)
		return
	}
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Println(err.Code.Class())
		}
		log.Println("error: apiuser.createUserAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	out.Nickname = in.Nickname
	out.Fullname = in.Fullname
	out.About = in.About
	out.Email = in.Email

	apiutil.WriteJsonObject(w, out, 201)
}
