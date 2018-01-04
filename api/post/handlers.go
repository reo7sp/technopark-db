package post

import (
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"net/http"
	"github.com/reo7sp/technopark-db/api/apimisc"
	"github.com/reo7sp/technopark-db/api/thread"
	"log"
	"strings"
	"github.com/reo7sp/technopark-db/api/user"
	"github.com/reo7sp/technopark-db/api/forum"
	"strconv"
)

func MakeShowDetailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		query := r.URL.Query()

		idStr := ps["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		relatedStr := query.Get("related")
		needUser := false
		needForum := false
		needThread := false
		for _, it := range strings.Split(relatedStr, ",") {
			switch it {
			case "user":
				needUser = true
			case "forum":
				needForum = true
			case "thread":
				needThread = true
			}
		}

		var postFullModel PostFullModel

		postModel, err := FindPostByIdInDB(db, id)
		if err != nil {
			errJson := apimisc.Error{Message: "Can't find post with id " + idStr}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}
		postFullModel.Post = postModel

		if needUser {
			userModel, err := user.FindUserByNicknameInDB(db, postModel.AuthorNickname)
			if err != nil {
				errJson := apimisc.Error{Message: "Can't find user with nickname " + postModel.AuthorNickname}
				apiutil.WriteJsonObject(w, errJson, 404)
				return
			}
			postFullModel.Author = userModel
		}

		if needForum {
			forumModel, err := forum.FindForumBySlugInDB(db, postModel.ForumSlug)
			if err != nil {
				errJson := apimisc.Error{Message: "Can't find forum with slug " + postModel.ForumSlug}
				apiutil.WriteJsonObject(w, errJson, 404)
				return
			}
			postFullModel.Forum = forumModel
		}

		if needThread {
			threadModel, err := thread.FindThreadByIdInDB(db, postModel.ThreadId)
			if err != nil {
				errJson := apimisc.Error{Message: "Can't find post with id " + strconv.FormatInt(postModel.ThreadId, 10)}
				apiutil.WriteJsonObject(w, errJson, 404)
				return
			}
			postFullModel.Thread = threadModel
		}

		apiutil.WriteJsonObject(w, postFullModel, 200)
	}
}

func MakeEditDetailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		idStr := ps["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		var postModel PostModel
		err = apiutil.ReadJsonObject(r, &postModel)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		rowsUpdated, err := EditMessageOfPostByIdInDB(id, postModel.Message, db)
		if err != nil {
			log.Println("error: api.post.MakeEditDetailsHandler: EditMessageOfPostByIdInDB:", err)
			w.WriteHeader(500)
			return
		}
		if rowsUpdated == 0 {
			errJson := apimisc.Error{Message: "Can't find post with id " + idStr}
			apiutil.WriteJsonObject(w, errJson, 404)
		}

		err = LoadUpPostHavingIdAndMessageFromDB(&postModel, db)
		if err != nil {
			log.Println("error: api.post.MakeEditDetailsHandler: LoadUpPostHavingIdAndMessageFromDB:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, postModel, 200)
	}
}
