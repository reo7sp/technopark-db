package apithread

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/patrickmn/go-cache"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MakeCreatePostHandler(db *pgx.ConnPool, cc *cache.Cache) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createPostRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createPostAction(w, in, db, cc)
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

func createPostCheckPostsParents(in createPostInput, db *pgx.ConnPool) (ok bool, err error) {
	parentIdsStrs := make([]string, 0, len(in.Posts))
	needToExecQuery := false
	for _, post := range in.Posts {
		if post.Parent == 0 {
			continue
		}
		needToExecQuery = true
		parentIdsStrs = append(parentIdsStrs, strconv.FormatInt(post.Parent, 10))
	}

	if !needToExecQuery {
		return true, nil
	}

	sqlQuery := fmt.Sprintf(`
	SELECT (
		CASE WHEN $1 IS TRUE
		THEN $2
		ELSE (SELECT id FROM threads WHERE slug = $3::citext)
		END
	) = ALL (SELECT threadId FROM posts WHERE id IN (%s))
	`, strings.Join(parentIdsStrs, ","))

	err = db.QueryRow(sqlQuery, in.HasId, in.Id, in.Slug).Scan(&ok)
	return
}

func createPostGenerateNextPlaceholder(i *int64) string {
	*i = *i + 1
	return "$" + strconv.FormatInt(*i, 10)
}

func createPostAction(w http.ResponseWriter, in createPostInput, db *pgx.ConnPool, cc *cache.Cache) {
	out := make(createPostOutput, 0, len(in.Posts))

	if len(in.Posts) == 0 {
		sqlQuery := "SELECT 1 FROM threads WHERE (CASE WHEN $1 IS TRUE THEN id = $2 ELSE slug = $3::citext END)"
		var ok int8
		err := db.QueryRow(sqlQuery, in.HasId, in.Id, in.Slug).Scan(&ok)

		if err != nil {
			errJson := api.ErrorModel{Message: "Can't find post thread"}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}

		apiutil.WriteJsonObject(w, out, 201)
		return
	}

	ok, err := createPostCheckPostsParents(in, db)
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

	sqlQuery := "INSERT INTO posts (parent, author, \"message\", forumSlug, threadId) VALUES"
	sqlValues := make([]interface{}, 0, 3+3*len(in.Posts))
	sqlValues = append(sqlValues, in.HasId, in.Id, in.Slug)
	placeholderIndex := int64(0 + 3)
	for i, post := range in.Posts {
		sqlQuery += fmt.Sprintf(` (
			%s,
			%s,
			%s,
			(SELECT forumSlug FROM threads WHERE (CASE WHEN $1 IS TRUE THEN id = $2 ELSE slug = $3::citext END)),
			(CASE WHEN $1 IS TRUE THEN $2 ELSE (SELECT id FROM threads WHERE slug = $3::citext) END)
		)`, createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex),
			createPostGenerateNextPlaceholder(&placeholderIndex))
		if i != len(in.Posts)-1 {
			sqlQuery += ","
		}
		sqlValues = append(sqlValues, post.ParentOrNil, post.Author, post.Message)
	}
	sqlQuery += " RETURNING id, createdAt, forumSlug::text, threadId"

	rows, err := db.Query(sqlQuery, sqlValues...)

	if err != nil {
		ok, constraint := dbutil.IsErrorAboutFailedForeignKeyReturnConstraint(err)

		if ok && constraint == "posts_parent_fkey" {
			errJson := api.ErrorModel{Message: "Parent post was created in another thread"}
			apiutil.WriteJsonObject(w, errJson, 409)
			return
		}
		if ok && constraint == "posts_author_fkey" {
			errJson := api.ErrorModel{Message: "Can't find post author"}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		}
		log.Println("error: apithread.createPostAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	i := 0
	for ; rows.Next(); i++ {
		var outItem createPostOutputItem

		var t time.Time
		err = rows.Scan(&outItem.Id, &t, &outItem.ForumSlug, &outItem.ThreadId)
		outItem.CreatedDateStr = t.UTC().Format(api.TIMEFORMAT)
		if err != nil {
			log.Println("error: apithread.createPostAction: INSERT scan iter:", err)
			w.WriteHeader(500)
			return
		}

		outItem.ParentPostId = &in.Posts[i].Parent
		outItem.AuthorNickname = in.Posts[i].Author
		outItem.Message = in.Posts[i].Message
		outItem.IsEdited = false

		out = append(out, outItem)
	}

	if i == 0 && len(in.Posts) != 0 {
		if in.Posts[0].Parent == 0 {
			errJson := api.ErrorModel{Message: "Can't find post author"}
			apiutil.WriteJsonObject(w, errJson, 404)
			return
		} else {
			errJson := api.ErrorModel{Message: "Parent post was created in another thread"}
			apiutil.WriteJsonObject(w, errJson, 409)
			return
		}
	}

	cc.IncrementInt64("posts_count", int64(len(in.Posts)))

	postsCount, ok := cc.Get("posts_count")
	if postsCount.(int64) == 1500000 {
		log.Println("info: apithread.createPostAction: run ANALYZE")
		db.Exec("ANALYZE")
	}

	apiutil.WriteJsonObject(w, out, 201)
}
