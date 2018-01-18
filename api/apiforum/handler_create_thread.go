package apiforum

import (
	"errors"
	"github.com/jackc/pgx"
	"github.com/patrickmn/go-cache"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"time"
)

func MakeCreateThreadHandler(db *pgx.ConnPool, cc *cache.Cache) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := createThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		createThreadAction(w, in, db, cc)
	}
	return f
}

type createThreadInput struct {
	ForumSlug string `json:"-"`

	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Message      string  `json:"message"`
	ThreadSlug   *string `json:"slug"`
	CreatedAtStr string  `json:"created"`
}

type createThreadOutput api.ThreadModel

func createThreadRead(r *http.Request, ps map[string]string) (in createThreadInput, err error) {
	slug, ok := ps["slug"]
	in.ForumSlug = slug
	if !ok {
		err = errors.New("slug is empty")
		return
	}

	err = apiutil.ReadJsonObject(r, &in)

	return
}

func createThreadAction(w http.ResponseWriter, in createThreadInput, db *pgx.ConnPool, cc *cache.Cache) {
	var out createThreadOutput

	if in.CreatedAtStr == "" {
		in.CreatedAtStr = time.Now().Format(time.RFC3339)
	}

	sqlQuery := "INSERT INTO threads (slug, title, author, forumSlug, \"message\", createdAt) VALUES ($1::citext, $2, $3::citext, (SELECT slug FROM forums WHERE slug = $4::citext), $5, $6) RETURNING id, forumSlug::text"
	err := db.QueryRow(sqlQuery, in.ThreadSlug, in.Title, in.Author, in.ForumSlug, in.Message, in.CreatedAtStr).Scan(&out.Id, &out.ForumSlug)

	if err != nil && dbutil.IsErrorAboutFailedForeignKey(err) {
		errJson := api.ErrorModel{Message: "Can't find user"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil && dbutil.IsErrorAboutDublicate(err) {
		sqlQuery := "SELECT id, title, author::text, \"message\", createdAt, slug::text, forumSlug::text, votesCount FROM threads WHERE slug = $1::citext"
		var t time.Time
		err := db.QueryRow(sqlQuery, in.ThreadSlug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.Message, &t, &out.Slug, &out.ForumSlug, &out.VotesCount)
		out.CreatedDateStr = t.UTC().Format(api.TIMEFORMAT)

		if err != nil {
			log.Println("error: apiforum.createThreadAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, out, 409)
		return
	}
	if err != nil {
		log.Println("error: apiforum.createThreadAction: INSERT:", err)
		w.WriteHeader(500)
		return
	}

	cc.IncrementInt64("threads_count", 1)

	out.Title = in.Title
	out.AuthorNickname = in.Author
	out.Message = in.Message
	out.CreatedDateStr = in.CreatedAtStr
	if in.ThreadSlug != nil {
		out.Slug = in.ThreadSlug
	}

	apiutil.WriteJsonObject(w, out, 201)
}
