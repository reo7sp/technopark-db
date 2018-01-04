package thread

import (
	"database/sql"
	"net/http"
	"strconv"
	"github.com/reo7sp/technopark-db/api/apimisc"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api/post"
)

type idOrSlugMatch struct {
	Id         int64
	Slug       string
	IsSlugAnId bool
}

func resolveIdOrSlugMatch(slug string) (m idOrSlugMatch) {
	m.Slug = slug

	if id, err := strconv.ParseInt(slug, 10, 64); err == nil {
		m.Id = id
		m.IsSlugAnId = true
	} else {
		m.IsSlugAnId = false
	}

	return
}

func make404JsonError(m idOrSlugMatch) apimisc.Error {
	if m.IsSlugAnId {
		return apimisc.Error{Message: "Can't find thread with id " + m.Slug}
	} else {
		return apimisc.Error{Message: "Can't find thread with slug " + m.Slug}
	}
}

func findThreadInDb(m idOrSlugMatch, db *sql.DB) (t ThreadModel, err error) {
	if m.IsSlugAnId {
		t, err = FindThreadByIdInDB(m.Id, db)
	} else {
		t, err = FindThreadBySlugInDB(m.Slug, db)
	}
	return
}

func MakeCreatePostHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		match := resolveIdOrSlugMatch(ps["slug"])

		var posts []post.PostModel
		err := apiutil.ReadJsonObject(r, &posts)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		// TODO

		apiutil.WriteJsonObject(w, posts, 200)
	}
}

func MakeShowDetailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		match := resolveIdOrSlugMatch(ps["slug"])

		threadModel, err := findThreadInDb(match, db)
		if err != nil {
			apiutil.WriteJsonObject(w, make404JsonError(match), 404)
			return
		}

		apiutil.WriteJsonObject(w, threadModel, 200)
	}
}

func MakeEditDetailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		match := resolveIdOrSlugMatch(ps["slug"])

		var threadModel ThreadModel
		err := apiutil.ReadJsonObject(r, &threadModel)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		threadModel, err = findThreadInDb(match, db)
		if err != nil {
			apiutil.WriteJsonObject(w, make404JsonError(match), 404)
			return
		}

		// TODO

		apiutil.WriteJsonObject(w, threadModel, 200)
	}
}

func MakeShowPostsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		match := resolveIdOrSlugMatch(ps["slug"])

		// TODO
	}
}

func MakeVoteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		match := resolveIdOrSlugMatch(ps["slug"])

		threadModel, err := findThreadInDb(match, db)
		if err != nil {
			apiutil.WriteJsonObject(w, make404JsonError(match), 404)
			return
		}

		var voteModel VoteModel
		err = apiutil.ReadJsonObject(r, &voteModel)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		threadModel.VotesCount += voteModel.Voice

		err = EditVotesCountOfThreadInDB(threadModel, db)
		if err != nil {
			log.Println("error: api.thread.MakeVoteHandler: EditVotesCountOfThreadInDB:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, threadModel, 200)
	}
}
