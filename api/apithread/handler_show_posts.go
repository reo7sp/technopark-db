package apithread

import (
	"net/http"
	"github.com/reo7sp/technopark-db/apiutil"
	"database/sql"
	"log"
	"github.com/reo7sp/technopark-db/api"
	"strconv"
)

func MakeShowPostsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
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
}

type showPostsOutputItem api.PostModel

type showPostsOutput []showPostsOutputItem

func showPostsRead(r *http.Request, ps map[string]string) (in showPostsInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)

	query := r.URL.Query()

	in.Limit, err = strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		err = nil
		in.Limit = -1
	}
	in.Since, err = strconv.ParseInt(query.Get("since"), 10, 64)
	if err != nil {
		err = nil
		in.Since = 0
	}
	in.Sort = query.Get("sort")
	in.IsDesc = query.Get("desc") == "true"

	return
}

func showPostsAction(w http.ResponseWriter, in showPostsInput, db *sql.DB) {
	var out showPostsOutput

	if in.Limit != -1 {
		out = make(showPostsOutput, 0, in.Limit)
	} else {
		out = make(showPostsOutput, 0)
	}

	sqlQuery := ""
	if in.Sort == "parent_tree" && in.IsDesc {
		sqlQuery += "WITH threadRootPostsCount AS (SELECT rootPostsCount FROM threads WHERE id = $1)"
	}
	sqlQuery += "SELECT id, parent, author, \"message\", isEdited, forumSlug, threadId, createdAt FROM posts"
	if in.HasId {
		sqlQuery += " WHERE (threadId = $1)"
	} else {
		sqlQuery += " WHERE (threadSlug = $1)"
	}
	sqlQuery += " AND (id >= $2)"
	switch in.Sort {
	case "flat":
		sqlQuery += " ORDER BY createdAt"
		if in.IsDesc {
			sqlQuery += " DESC"
		} else {
			sqlQuery += " ASC"
		}
		sqlQuery += " LIMIT $3"

	case "tree":
		sqlQuery += " ORDER BY path"
		if in.IsDesc {
			sqlQuery += " DESC"
		} else {
			sqlQuery += " ASC"
		}
		sqlQuery += " LIMIT $3"

	case "parent_tree":
		if in.IsDesc {
			sqlQuery += " AND (rootPostNo < $3)"
		} else {
			sqlQuery += " AND (rootPostNo >= threadRootPostsCount - $3)"
		}
		sqlQuery += " ORDER BY path"
		if in.IsDesc {
			sqlQuery += " DESC"
		} else {
			sqlQuery += " ASC"
		}
	}

	log.Println(sqlQuery, in.Slug, in.Since, in.Limit)
	rows, err := db.Query(sqlQuery, in.Slug, in.Since, in.Limit)
	if err != nil {
		log.Println("error: apithread.showPostsAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showPostsOutputItem
		err = rows.Scan(
			&outItem.Id, &outItem.ParentPostId, &outItem.AuthorNickname, &outItem.Message, &outItem.IsEdited,
			&outItem.ForumSlug, &outItem.ThreadId, &outItem.CreatedDateStr)
		if err != nil {
			log.Println("error: apithread.showPostsAction: SELECT iter:", err)
			w.WriteHeader(500)
			return
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
