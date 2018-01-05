package apipost

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"strconv"
	"github.com/reo7sp/technopark-db/api"
)

func MakeEditPostHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := editPostRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		editPostAction(w, in, db)
	}
	return f
}

type editPostInput struct {
	Id      int64  `json:"-"`
	IdStr   string `json:"-"`
	Message string `json:"message"`
}

type editPostOutput api.PostModel

func editPostRead(r *http.Request, ps map[string]string) (in editPostInput, err error) {
	in.IdStr = ps["id"]
	id, err := strconv.ParseInt(in.IdStr, 10, 64)
	if err != nil {
		return
	}
	in.Id = id

	err = apiutil.ReadJsonObject(r, &in)

	return
}

func editPostAction(w http.ResponseWriter, in editPostInput, db *sql.DB) {
	var out editPostOutput

	sqlQuery := "UPDATE posts SET \"message\" = $1, isEdited = TRUE WHERE id = $2"
	r, err := db.Exec(sqlQuery, in.Message, in.Id)
	if err != nil {
		log.Println("error: apipost.editPostAction: UPDATE:", err)
		w.WriteHeader(500)
		return
	}
	if c, _ := r.RowsAffected(); c == 0 {
		errJson := api.Error{Message: "Can't find post with id " + in.IdStr}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
