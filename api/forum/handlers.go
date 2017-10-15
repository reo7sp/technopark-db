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
		var forum ForumModel
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
					return
				}
			} else {
				w.WriteHeader(500)
				log.Println("error: api.CreateFuncMaker: forum.Create:", err)
				return
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

func DetailsFuncMaker(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		slug := ps.ByName("slug")
		forum, err := FindForum(db, slug)
		if err != nil {
			apiErr := api.Error{Message: "Can't find forum with slug " + slug}
			api.WriteJsonObject(w, apiErr, 404)
			return
		}

		api.WriteJsonObject(w, forum, 200)
	}
}
