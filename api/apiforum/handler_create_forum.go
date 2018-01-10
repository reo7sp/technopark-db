package apiforum

import (
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"github.com/jackc/pgx"
)

func MakeCreateForumHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createForumRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createForumAction(w, in, db)
	}
	return f
}

type createForumInput struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type createForumOutput api.ForumModel

type createForumGetUserInfo struct {
	Nickname string
}

func createForumRead(r *http.Request, ps map[string]string) (in createForumInput, err error) {
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func createForumGetUser(in createForumInput, db *pgx.ConnPool) (r createForumGetUserInfo, err error) {
	sqlQuery := "SELECT nickname FROM users WHERE nickname = $1"
	err = db.QueryRow(sqlQuery, in.User).Scan(&r.Nickname)
	return
}

func createForumAction(w http.ResponseWriter, in createForumInput, db *pgx.ConnPool) {
	forumInfo, err := createForumGetUser(in, db)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find user"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	var out createForumOutput

	sqlQuery := "INSERT INTO forums (title, \"user\", slug) VALUES ($1, $2, $3)"
	_, err = db.Exec(sqlQuery, in.Title, forumInfo.Nickname, in.Slug)

	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT slug, title, \"user\", postsCount, threadsCount FROM forums WHERE slug = $1"
		err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Slug, &out.Title, &out.User, &out.PostsCount, &out.ThreadsCount)

		if err != nil {
			log.Println("error: apiforum.createForumAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, out, 409)
		return
	}
	if err != nil {
		log.Println("error: apiforum.createForumAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	out.Slug = in.Slug
	out.Title = in.Title
	out.User = forumInfo.Nickname
	out.PostsCount = 0
	out.ThreadsCount = 0

	apiutil.WriteJsonObject(w, out, 201)
}
