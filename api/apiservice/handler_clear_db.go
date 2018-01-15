package apiservice

import (
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"github.com/patrickmn/go-cache"
)

func MakeClearDbHandler(db *pgx.ConnPool, cc *cache.Cache) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		clearDbAction(w, db, cc)
	}
	return f
}

func clearDbAction(w http.ResponseWriter, db *pgx.ConnPool, cc *cache.Cache) {
	_, err := db.Exec("TRUNCATE TABLE forums, threads, users, posts RESTART IDENTITY CASCADE")
	if err != nil {
		log.Println("error: apiservice.clearDbAction: TRUNCATE:", err)
		w.WriteHeader(500)
		return
	}

	cc.Set("forums_count", int64(0), cache.NoExpiration)
	cc.Set("threads_count", int64(0), cache.NoExpiration)
	cc.Set("users_count", int64(0), cache.NoExpiration)
	cc.Set("posts_count", int64(0), cache.NoExpiration)

	w.WriteHeader(200)
}
