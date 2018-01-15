package apiservice

import (
	"github.com/jackc/pgx"
	"github.com/patrickmn/go-cache"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"net/http"
)

func MakeShowStatusHandler(db *pgx.ConnPool, cc *cache.Cache) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		showStatusAction(w, db, cc)
	}
	return f
}

type showStatusOutput api.StatusModel

func showStatusAction(w http.ResponseWriter, db *pgx.ConnPool, cc *cache.Cache) {
	var out showStatusOutput

	forumsCount, ok1 := cc.Get("forums_count")
	threadsCount, ok2 := cc.Get("threads_count")
	usersCount, ok3 := cc.Get("users_count")
	postsCount, ok4 := cc.Get("posts_count")

	if ok1 && ok2 && ok3 && ok4 {
		out.ForumsCount = forumsCount.(int64)
		out.ThreadsCount = threadsCount.(int64)
		out.UsersCount = usersCount.(int64)
		out.PostsCount = postsCount.(int64)
	} else {
		var err error
		out.ForumsCount, out.ThreadsCount, out.UsersCount, out.PostsCount, err = LoadCountsToCache(db, cc)
		if err != nil {
			log.Println("error: apiservice.showStatusAction: SELECT:", err)
			w.WriteHeader(500)
			return
		}
	}

	apiutil.WriteJsonObject(w, out, 200)
}

func LoadCountsToCache(db *pgx.ConnPool, cc *cache.Cache) (f, t, u, p int64, err error) {
	sqlQuery := "SELECT (SELECT count(*) FROM forums), (SELECT count(*) FROM threads), (SELECT count(*) FROM users), (SELECT count(*) FROM posts)"
	err = db.QueryRow(sqlQuery).Scan(&f, &t, &u, &p)
	if err != nil {
		return
	}

	cc.Set("forums_count", f, cache.NoExpiration)
	cc.Set("threads_count", t, cache.NoExpiration)
	cc.Set("users_count", u, cache.NoExpiration)
	cc.Set("posts_count", p, cache.NoExpiration)
	return
}
