package forum

import (
	"database/sql"
	"github.com/reo7sp/technopark-db/api/thread"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/api/apimisc"
	"github.com/reo7sp/technopark-db/api/user"
)

func MakeCreateForumHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		var forumModel ForumModel
		err := apiutil.ReadJsonObject(r, &forumModel)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		err = CreateForumInDB(forumModel, db)
		isDublicate := false
		if err != nil {
			if dbutil.IsErrorAboutDublicate(err) {
				isDublicate = true
			} else {
				log.Println("error: api.forum.MakeCreateForumHandler: CreateForumInDB:", err)
				w.WriteHeader(500)
				return
			}
		}

		if isDublicate {
			forumModel, err = FindForumBySlugInDB(forumModel.Slug, db)
		}

		statusCode := 201
		if isDublicate {
			statusCode = 409
		}
		apiutil.WriteJsonObject(w, forumModel, statusCode)
	}
}

func MakeShowDetailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		slug := ps["slug"]

		forumModel, err := FindForumBySlugInDB(slug, db)
		if err != nil {
			errJson := apimisc.Error{Message: "Can't find forum with slug " + slug}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}

		apiutil.WriteJsonObject(w, forumModel, 200)
	}
}

func MakeCreateThreadHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		slug := ps["slug"]

		var threadModel thread.ThreadModel
		err := apiutil.ReadJsonObject(r, &threadModel)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		threadModel.Slug = slug

		err = thread.CreateThreadInDB(threadModel, db)
		isDublicate := false
		if err != nil {
			if dbutil.IsErrorAboutDublicate(err) {
				isDublicate = true
				threadModel, err = thread.FindThreadBySlugInDB(slug, db)
			}
			if err != nil {
				w.WriteHeader(500)
				log.Println("error: api.forum.MakeCreateForumHandler: threadModel.Create:", err)
				return
			}
		}

		statusCode := 201
		if isDublicate {
			statusCode = 409
		}
		apiutil.WriteJsonObject(w, threadModel, statusCode)
	}
}

func MakeShowUsersHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		query := r.URL.Query()

		slug := ps["slug"]
		limit := query.Get("limit")
		since := query.Get("since")
		isDesc := query.Get("desc")

		users, err := user.FindUsersByForumInDB(db, slug, limit, since, isDesc)
		if err != nil {
			errJson := apimisc.Error{Message: "Can't find forum with slug " + slug}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}

		apiutil.WriteJsonObject(w, users, 200)
	}
}

func MakeShowThreadsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		query := r.URL.Query()

		slug := ps["slug"]
		limit := query.Get("limit")
		since := query.Get("since")
		isDesc := query.Get("desc")

		threads, err := thread.FindThreadsByForumInDB(db, slug, limit, since, isDesc)
		if err != nil {
			errJson := apimisc.Error{Message: "Can't find forum with slug " + slug}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}

		apiutil.WriteJsonObject(w, threads, 200)
	}
}
