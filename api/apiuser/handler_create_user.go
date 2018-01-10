package apiuser

import (
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"github.com/jackc/pgx"
)

func MakeCreateUserHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
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

type createUserOutput2 []api.UserModel

func createUserRead(r *http.Request, ps map[string]string) (in createUserInput, err error) {
	in.Nickname = ps["nickname"]
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func createUserAction(w http.ResponseWriter, in createUserInput, db *pgx.ConnPool) {
	var out createUserOutput

	sqlQuery := "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlQuery, in.Nickname, in.Fullname, in.About, in.Email)

	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		out2 := make(createUserOutput2, 0, 2)

		sqlQuery := "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2"
		rows, err := db.Query(sqlQuery, in.Nickname, in.Email)

		if err != nil {
			log.Println("error: apiforum.createForumAction: SELECT by nickname start:", err)
			w.WriteHeader(500)
			return
		}

		defer rows.Close()
		for rows.Next() {
			var user api.UserModel

			err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
			if err != nil {
				log.Println("error: apiforum.createForumAction: SELECT by nickname iter:", err)
				w.WriteHeader(500)
				return
			}

			out2 = append(out2, user)
		}

		apiutil.WriteJsonObject(w, out2, 409)
		return
	}
	if err != nil {
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
