package apiuser

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
)

func MakeEditUserHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := editUserRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		editUserAction(w, in, db)
	}
	return f
}

type editUserInput struct {
	Nickname string `json:"-"`

	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type editUserOutput api.UserModel

func editUserRead(r *http.Request, ps map[string]string) (in editUserInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func editUserAction(w http.ResponseWriter, in editUserInput, db *sql.DB) {
	var out editUserOutput

	sqlQuery := "UPDATE users SET fullname = $1, about = $2, email = $3 WHERE nickname = $4"
	_, err := db.Exec(sqlQuery, in.Fullname, in.About, in.Email, in.Nickname)

	if err != nil {
		log.Println("error: apiuser.editUserAction: UPDATE:", err)
		w.WriteHeader(500)
		return
	}

	out.Nickname = in.Nickname
	out.Fullname = in.Fullname
	out.About = in.About
	out.Email = in.Email

	apiutil.WriteJsonObject(w, out, 200)
}
