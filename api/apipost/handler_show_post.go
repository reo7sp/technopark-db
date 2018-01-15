package apipost

import (
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MakeShowPostHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showPostRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showPostAction(w, in, db)
	}
	return f
}

type showPostInput struct {
	Id         int64
	NeedUser   bool
	NeedForum  bool
	NeedThread bool
}

type showPostOutputBuilder struct {
	Post   api.PostModel
	Author api.UserModel
	Thread api.ThreadModel
	Forum  api.ForumModel
}

type showPostOutput map[string]interface{}

func showPostRead(r *http.Request, ps map[string]string) (in showPostInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)

	query := r.URL.Query()

	idStr := ps["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}
	in.Id = id

	relatedStr := query.Get("related")
	for _, it := range strings.Split(relatedStr, ",") {
		switch it {
		case "user":
			in.NeedUser = true
		case "forum":
			in.NeedForum = true
		case "thread":
			in.NeedThread = true
		}
	}

	return
}

func showPostAction(w http.ResponseWriter, in showPostInput, db *pgx.ConnPool) {
	var outBuilder showPostOutputBuilder
	out := make(showPostOutput)

	const maxCountOfSqlScans = 24

	sqlJoins := ""
	sqlFields := ""
	sqlScans := make([]interface{}, 0, maxCountOfSqlScans)

	sqlFields += " p.id, p.parent, p.author::text, p.message, p.isEdited, p.forumSlug::text, p.threadId, p.createdAt"
	var t1 time.Time
	sqlScans = append(sqlScans,
		&outBuilder.Post.Id, &outBuilder.Post.ParentPostId, &outBuilder.Post.AuthorNickname, &outBuilder.Post.Message,
		&outBuilder.Post.IsEdited, &outBuilder.Post.ForumSlug, &outBuilder.Post.ThreadId, &t1)

	if in.NeedUser {
		sqlJoins += " JOIN users u ON (u.nickname = p.author)"
		sqlFields += ", u.nickname::text, u.fullname, u.email::text, u.about"
		sqlScans = append(sqlScans,
			&outBuilder.Author.Nickname, &outBuilder.Author.Fullname, &outBuilder.Author.Email, &outBuilder.Author.About)
	}

	if in.NeedForum {
		sqlJoins += " JOIN forums f ON (f.slug = p.forumSlug)"
		sqlFields += ", f.slug::text, f.title, f.user::text, f.postsCount, f.threadsCount"
		sqlScans = append(sqlScans,
			&outBuilder.Forum.Slug, &outBuilder.Forum.Title, &outBuilder.Forum.User,
			&outBuilder.Forum.PostsCount, &outBuilder.Forum.ThreadsCount)
	}

	var t2 time.Time
	if in.NeedThread {
		sqlJoins += " JOIN threads t ON (t.id = p.threadId)"
		sqlFields += ", t.id, t.title, t.author::text, t.forumSlug::text, t.message, t.createdAt, t.slug::text, t.votesCount"
		sqlScans = append(sqlScans,
			&outBuilder.Thread.Id, &outBuilder.Thread.Title, &outBuilder.Thread.AuthorNickname, &outBuilder.Thread.ForumSlug,
			&outBuilder.Thread.Message, &t2, &outBuilder.Thread.Slug, &outBuilder.Thread.VotesCount)
	}

	sqlQuery := "SELECT " + sqlFields + " FROM posts p " + sqlJoins + " WHERE p.id = $1"

	err := db.QueryRow(sqlQuery, in.Id).Scan(sqlScans...)

	outBuilder.Post.CreatedDateStr = t1.UTC().Format(api.TIMEFORMAT)
	if in.NeedThread {
		outBuilder.Thread.CreatedDateStr = t2.UTC().Format(api.TIMEFORMAT)
	}

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find post"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apipost.showPostAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	out["post"] = outBuilder.Post
	if in.NeedUser {
		out["author"] = outBuilder.Author
	}
	if in.NeedForum {
		out["forum"] = outBuilder.Forum
	}
	if in.NeedThread {
		out["thread"] = outBuilder.Thread
	}

	apiutil.WriteJsonObject(w, out, 200)
}
