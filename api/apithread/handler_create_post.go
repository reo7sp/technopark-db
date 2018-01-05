package apithread

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api/apipost"
	"errors"
)

func MakeCreatePostHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createPostRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createPostAction(w, in, db)
	}
	return f
}

type createPostInputItem struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`

	ParentOrNil interface{} `json:"-"`
}

type createPostInput struct {
	slugOrIdInput

	Posts []createPostInputItem
}

type createPostOutputItem apipost.PostModel

type createPostOutput []createPostOutputItem

type createPostGetThreadInfo struct {
	Id        int64
	Slug      string
	ForumSlug string
}

func createPostRead(r *http.Request, ps map[string]string) (in createPostInput, err error) {
	resolveSlugOrIdInput(ps["slug"], &in.slugOrIdInput)

	err = apiutil.ReadJsonObject(r, &in.Posts)
	if err != nil {
		return
	}

	if len(in.Posts) == 0 {
		err = errors.New("posts in empty")
		return
	}

	for _, post := range in.Posts {
		if post.Parent == 0 {
			post.ParentOrNil = nil
		} else {
			post.ParentOrNil = post.Parent
		}
	}

	return
}

func createPostGetThread(in createPostInput, db *sql.DB) (r createPostGetThreadInfo, err error) {
	if in.HasId {
		r.Id = in.Id
		sqlQuery := "SELECT slug, forumSlug FROM threads WHERE id = $1"
		err = db.QueryRow(sqlQuery, in.Id).Scan(&r.Slug, &r.ForumSlug)
	} else {
		sqlQuery := "SELECT id, forumSlug FROM threads WHERE slug = $1"
		err = db.QueryRow(sqlQuery, in.Slug).Scan(&r.Id, &r.ForumSlug)
	}
	return
}

func createPostAction(w http.ResponseWriter, in createPostInput, db *sql.DB) {
	threadInfo, err := createPostGetThread(in, db)
	if err != nil {
		log.Println("error: apithread.createPostAction: createPostGetThread:", err)
		w.WriteHeader(500)
		return
	}

	out := make(createPostOutput, 0, len(in.Posts))

	sqlQuery := "INSERT INTO posts (parent, author, \"message\", forumSlug, threadId, threadSlug) VALUES"
	sqlValues := make([]interface{}, 0, 5*len(in.Posts))
	for i, post := range in.Posts {
		sqlQuery += " (?, ?, ?, ?, ?)"
		if i != len(in.Posts)-1 {
			sqlQuery += ","
		}
		sqlValues = append(sqlValues, post.ParentOrNil, post.Author, post.Message, threadInfo.ForumSlug, threadInfo.Id, threadInfo.Slug)
	}
	sqlQuery += " RETURNING id, createdAt"

	rows, err := db.Query(sqlQuery, sqlValues)
	if err != nil {
		log.Println("error: apithread.createPostAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		var outItem createPostOutputItem
		err = rows.Scan(&outItem.Id, &outItem.CreatedDateStr)
		if err != nil {
			log.Println("error: apithread.createPostAction: INSERT scan iter:", err)
			w.WriteHeader(500)
			return
		}
		outItem.ParentPostId = in.Posts[i].Parent
		outItem.AuthorNickname = in.Posts[i].Author
		outItem.Message = in.Posts[i].Message
		outItem.IsEdited = false
		outItem.ForumSlug = threadInfo.ForumSlug
		outItem.ThreadId = threadInfo.Id
	}

	apiutil.WriteJsonObject(w, out, 200)
}
