package forum

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/database"
	"github.com/reo7sp/technopark-db/api"
)

func CreateFuncMaker(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var forum Forum
		err := api.ReadJsonObject(w, r, &forum)
		if err != nil {
			return
		}

		err = forum.Create(db)
		isDublicate := false
		if err != nil {
			if database.IsErrorAboutDublicate(err) {
				isDublicate = true
				err = forum.FetchPostsAndThreadsCount(db)
				if err != nil {
					w.WriteHeader(500)
					log.Println("error: api.CreateFuncMaker: forum.FetchPostsAndThreadsCount:", err)
				}
			} else {
				w.WriteHeader(500)
				log.Println("error: api.CreateFuncMaker: forum.Create:", err)
			}
		}

		var statusCode int
		if isDublicate {
			statusCode = 201
		} else {
			statusCode = 409
		}
		api.WriteJsonObject(w, forum, statusCode)
	}
}
