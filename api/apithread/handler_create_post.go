package apithread

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api"
	"strconv"
	"fmt"
	"github.com/reo7sp/technopark-db/dbutil"
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

	ParentOrNil *int64 `json:"-"`
}

type createPostInput struct {
	slugOrIdInput

	Posts []createPostInputItem
}

type createPostOutputItem api.PostModel

type createPostOutput []createPostOutputItem

type createPostGetThreadInfo struct {
	Id        int64
	Slug      *string
	ForumSlug string
}

func createPostRead(r *http.Request, ps map[string]string) (in createPostInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)

	err = apiutil.ReadJsonObject(r, &in.Posts)
	if err != nil {
		return
	}

	if len(in.Posts) == 0 {
		return
	}

	for i, post := range in.Posts {
		if post.Parent == 0 {
			post.ParentOrNil = nil
		} else {
			in.Posts[i].ParentOrNil = &in.Posts[i].Parent
		}
	}

	return
}

func createPostCheckPostsParents(threadId int64, Posts []createPostInputItem, db *sql.DB) (bool, error) {
	sqlQuery := "SELECT threadId FROM posts WHERE id IN (0"
	needToExecQuery := false
	for _, post := range Posts {
		if post.Parent == 0 {
			continue
		}
		needToExecQuery = true
		sqlQuery += ", " + strconv.FormatInt(post.Parent, 10)
	}
	sqlQuery += ")"

	if !needToExecQuery {
		return true, nil
	}

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		if id != threadId {
			return false, nil
		}
	}
	return true, nil
}

func createPostGetThread(in createPostInput, db *sql.DB) (r createPostGetThreadInfo, err error) {
	if in.HasId {
		r.Id = in.Id
		sqlQuery := "SELECT slug, forumSlug FROM threads WHERE id = $1"
		err = db.QueryRow(sqlQuery, in.Id).Scan(&r.Slug, &r.ForumSlug)
	} else {
		r.Slug = &in.Slug
		sqlQuery := "SELECT id, forumSlug FROM threads WHERE slug = $1"
		err = db.QueryRow(sqlQuery, in.Slug).Scan(&r.Id, &r.ForumSlug)
	}
	return
}

func createPostGenerateNextPlaceholder(i *int64) string {
	*i = *i + 1
	return "$" + strconv.FormatInt(*i, 10)
}

func createPostAction(w http.ResponseWriter, in createPostInput, db *sql.DB) {
	out := make(createPostOutput, 0, len(in.Posts))

	threadInfo, err := createPostGetThread(in, db)
	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find post thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.createPostAction: createPostGetThread:", err)
		w.WriteHeader(500)
		return
	}

	if len(in.Posts) == 0 {
		apiutil.WriteJsonObject(w, out, 201)
		return
	}

	ok, err := createPostCheckPostsParents(threadInfo.Id, in.Posts, db)
	if err != nil {
		log.Println("error: apithread.createPostAction: createPostCheckPostsParents:", err)
		w.WriteHeader(500)
		return
	}
	if !ok {
		errJson := api.ErrorModel{Message: "Parent post was created in another thread"}
		apiutil.WriteJsonObject(w, errJson, 409)
		return
	}


	sqlQuery := "INSERT INTO posts (parent, author, \"message\", forumSlug, threadId, threadSlug) VALUES"
	sqlValues := make([]interface{}, 0, 6*len(in.Posts))
	placeholderIndex := int64(0)
	for i, post := range in.Posts {
		sqlQuery += fmt.Sprintf(" (%s, %s, %s, %s, %s, %s)",
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex))
		if i != len(in.Posts)-1 {
			sqlQuery += ","
		}
		sqlValues = append(sqlValues, post.ParentOrNil, post.Author, post.Message, threadInfo.ForumSlug, threadInfo.Id, threadInfo.Slug)
	}
	sqlQuery += " RETURNING id, createdAt"

	rows, err := db.Query(sqlQuery, sqlValues...)

	if err != nil {
		ok, constaint := dbutil.IsErrorAboutFailedForeignKeyReturnConstaint(err)

		if ok && constaint == "posts_parent_fkey" {
			errJson := api.ErrorModel{Message: "Parent post was created in another thread"}
			apiutil.WriteJsonObject(w, errJson, 409)
			return
		}
		if ok && constaint == "posts_author_fkey" {
			errJson := api.ErrorModel{Message: "Can't find post author"}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}
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

		outItem.ParentPostId = &in.Posts[i].Parent
		outItem.AuthorNickname = in.Posts[i].Author
		outItem.Message = in.Posts[i].Message
		outItem.IsEdited = false
		outItem.ForumSlug = threadInfo.ForumSlug
		outItem.ThreadId = threadInfo.Id

		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 201)
}
