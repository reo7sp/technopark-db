package apithread

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func MakeShowPostsHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showPostsRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showPostsAction(w, in, db)
	}
	return f
}

type showPostsInput struct {
	slugOrIdInput

	Since  int64  `json:"since"`
	Sort   string `json:"sort"`
	Limit  int64  `json:"limit"`
	IsDesc bool   `json:"desc"`

	Order string `json:"-"`
}

type showPostsOutputItem api.PostModel

type showPostsOutput []showPostsOutputItem

func showPostsRead(r *http.Request, ps map[string]string) (in showPostsInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)

	query := r.URL.Query()

	in.Limit, err = strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		return
	}
	in.Since, err = strconv.ParseInt(query.Get("since"), 10, 64)
	if err != nil {
		err = nil
		in.Since = -1
	}
	in.Sort = query.Get("sort")
	if in.Sort == "" {
		in.Sort = "flat"
	}
	in.IsDesc = query.Get("desc") == "true"
	if in.IsDesc {
		in.Order = "DESC"
	} else {
		in.Order = "ASC"
	}

	return
}

func showPostsCheckThreadExists(in slugOrIdInput, db *pgx.ConnPool) (bool, error) {
	sqlQuery := "SELECT id FROM threads WHERE"
	if in.HasId {
		sqlQuery += " id = $1"
	} else {
		sqlQuery += " slug = $1"
	}
	var i int64
	err := db.QueryRow(sqlQuery, in.Slug).Scan(&i)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func showPostsAction(w http.ResponseWriter, in showPostsInput, db *pgx.ConnPool) {
	doesThreadExists, err := showPostsCheckThreadExists(in.slugOrIdInput, db)
	if err != nil {
		log.Println("error: apithread.showPostsAction: showPostsCheckThreadExists:", err)
		w.WriteHeader(500)
		return
	}
	if !doesThreadExists {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	var out showPostsOutput

	if in.Limit != -1 {
		out = make(showPostsOutput, 0, in.Limit)
	} else {
		out = make(showPostsOutput, 0)
	}

	var sqlQuery string
	var sqlValues []interface{}

	switch in.Sort {
	case "flat":
		sqlQuery = fmt.Sprintf(`

		SELECT id, parent, author::text, "message", isEdited, forumSlug::text, threadId, createdAt FROM posts
		WHERE threadId = (
			CASE WHEN $1 IS TRUE
			THEN $2
			ELSE (SELECT id FROM threads WHERE slug = $3::citext)
			END
		)
		AND (
			CASE WHEN $4 != -1
			THEN (
				CASE WHEN $6 IS TRUE
				THEN (id < $4)
				ELSE (id > $4)
				END
			)
			ELSE TRUE
			END
		)
		ORDER BY createdAt %s, id %s
		LIMIT $5::bigint

		`, in.Order, in.Order)

		sqlValues = []interface{}{in.HasId, in.Id, in.Slug, in.Since, in.Limit, in.IsDesc}

	case "tree":
		sqlQuery = fmt.Sprintf(`

		SELECT id, parent, author::text, "message", isEdited, forumSlug::text, threadId, createdAt FROM posts
		WHERE threadId = (
			CASE WHEN $1 IS TRUE
			THEN $2
			ELSE (SELECT id FROM threads WHERE slug = $3::citext)
			END
		)
		AND (
			CASE WHEN $4 != -1
			THEN (
				CASE WHEN $6 IS TRUE
				THEN (path < (SELECT p1.path FROM posts p1 WHERE p1.id = $4))
				ELSE (path > (SELECT p1.path FROM posts p1 WHERE p1.id = $4))
				END
			)
			ELSE TRUE
			END
		)
		ORDER BY path %s, id %s
		LIMIT $5::bigint

		`, in.Order, in.Order)

		sqlValues = []interface{}{in.HasId, in.Id, in.Slug, in.Since, in.Limit, in.IsDesc}

	case "parent_tree":
		sqlQuery = fmt.Sprintf(`

		WITH rootParentPosts AS (
			SELECT id FROM posts
			WHERE threadId = (
				CASE WHEN $1 IS TRUE
				THEN $2
				ELSE (SELECT id FROM threads WHERE slug = $3::citext)
				END
			)
			AND (
				CASE WHEN $4 != -1
				THEN (
					CASE WHEN $6 IS TRUE
					THEN (path < (SELECT path FROM posts p1 WHERE p1.id = $4))
					ELSE (path > (SELECT path FROM posts p1 WHERE p1.id = $4))
					END
				)
				ELSE TRUE
				END
			)
			AND parent IS NULL
			ORDER BY path %s, id %s
			LIMIT $5
		)
		SELECT p.id, p.parent, p.author::text, p.message, p.isEdited, p.forumSlug::text, p.threadId, p.createdAt
		FROM rootParentPosts rpp
		JOIN posts p ON (p.rootparent = rpp.id)
		ORDER BY path %s, id %s

		`, in.Order, in.Order, in.Order, in.Order)

		sqlValues = []interface{}{in.HasId, in.Id, in.Slug, in.Since, in.Limit, in.IsDesc}
	}

	rows, err := db.Query(sqlQuery, sqlValues...)
	if err != nil {
		log.Println("error: apithread.showPostsAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showPostsOutputItem
		var t time.Time
		err = rows.Scan(
			&outItem.Id, &outItem.ParentPostId, &outItem.AuthorNickname, &outItem.Message, &outItem.IsEdited,
			&outItem.ForumSlug, &outItem.ThreadId, &t)
		outItem.CreatedDateStr = t.UTC().Format(api.TIMEFORMAT)
		if err != nil {
			log.Println("error: apithread.showPostsAction: SELECT iter:", err)
			w.WriteHeader(500)
			return
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
