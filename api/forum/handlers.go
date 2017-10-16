package forum

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/database"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/api/thread"
)

func CreateFuncMaker(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var forumModel ForumModel
		err := api.ReadJsonObject(w, r, &forumModel)
		if err != nil {
			return
		}

		err = forumModel.Create(db)
		var isDublicate bool
		if err != nil {
			if database.IsErrorAboutDublicate(err) {
				isDublicate = true
				forumModel, err = FindForum(db, forumModel.Slug)
			}
			if err != nil {
				w.WriteHeader(500)
				log.Println("error: api.CreateFuncMaker: forumModel.Create:", err)
				return
			}
		} else {
			isDublicate = false
		}

		var statusCode int
		if isDublicate {
			statusCode = 201
		} else {
			statusCode = 409
		}
		api.WriteJsonObject(w, forumModel, statusCode)
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

func CreateThreadFuncMaker(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		slug := ps.ByName("slug")

		var threadModel thread.ThreadModel
		err := api.ReadJsonObject(w, r, &threadModel)
		threadModel.Slug = slug
		if err != nil {
			return
		}

		err = threadModel.Create(db)
		var isDublicate bool
		if err != nil {
			if database.IsErrorAboutDublicate(err) {
				isDublicate = true
				threadModel, err = thread.FindThread(db, slug)
			}
			if err != nil {
				w.WriteHeader(500)
				log.Println("error: api.CreateFuncMaker: threadModel.Create:", err)
				return
			}
		} else {
			isDublicate = false
		}

		var statusCode int
		if isDublicate {
			statusCode = 201
		} else {
			statusCode = 409
		}
		api.WriteJsonObject(w, threadModel, statusCode)
	}
}
