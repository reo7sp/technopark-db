package apiuser

import (
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
	"github.com/jackc/pgx"
)

func MakeEditUserHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
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
	in.Nickname = ps["nickname"]
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func editUserAction(w http.ResponseWriter, in editUserInput, db *pgx.ConnPool) {
	var out editUserOutput

	sqlQuery := "UPDATE users SET"
	sqlValues := make([]interface{}, 0, 4)
	if in.Fullname != "" {
		sqlValues = append(sqlValues, in.Fullname)
		sqlQuery += " fullname = $1"
	}
	if in.About != "" {
		if len(sqlValues) != 0 {
			sqlQuery += ","
		}
		sqlValues = append(sqlValues, in.About)
		sqlQuery += " about = $" + strconv.FormatInt(int64(len(sqlValues)), 10)
	}
	if in.Email != "" {
		if len(sqlValues) != 0 {
			sqlQuery += ","
		}
		sqlValues = append(sqlValues, in.Email)
		sqlQuery += " email = $" + strconv.FormatInt(int64(len(sqlValues)), 10)
	}
	if len(sqlValues) != 0 {
		sqlValues = append(sqlValues, in.Nickname)
		sqlQuery += " WHERE nickname = $" + strconv.FormatInt(int64(len(sqlValues)), 10)
	}

	if len(sqlValues) != 0 {
		_, err := db.Exec(sqlQuery, sqlValues...)

		if err != nil && dbutil.IsErrorAboutDublicate(err) {
			errJson := api.ErrorModel{Message: "This email is already registered by user"}
			apiutil.WriteJsonObject(w, errJson, 409)
			return
		}
		if err != nil {
			log.Println("error: apiuser.editUserAction: UPDATE:", err)
			w.WriteHeader(500)
			return
		}
	}

	sqlQuery = "SELECT nickname::text, fullname, about, email::text FROM users WHERE nickname = $1::citext"
	err := db.QueryRow(sqlQuery, in.Nickname).Scan(&out.Nickname, &out.Fullname, &out.About, &out.Email)
	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find user"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apiuser.editUserAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
