package main

import (
	"context"
	"github.com/dimfeld/httptreemux"
	"github.com/jackc/pgx"
	"github.com/patrickmn/go-cache"
	"github.com/reo7sp/technopark-db/api/apiforum"
	"github.com/reo7sp/technopark-db/api/apipost"
	"github.com/reo7sp/technopark-db/api/apiservice"
	"github.com/reo7sp/technopark-db/api/apithread"
	"github.com/reo7sp/technopark-db/api/apiuser"
	"github.com/reo7sp/technopark-db/dbutil"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Connecting to db")
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}

	cc := cache.New(cache.NoExpiration, cache.NoExpiration)

	log.Println("Loading counts to cache")
	apiservice.LoadCountsToCache(db, cc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)

	web := setupHttpServer(db, cc)
	go func() {
		log.Println("Starting http server: http://localhost:5000")
		if err := web.ListenAndServe(); err != nil {
			log.Println("Web stop:", err)
		}
	}()
	<-stop

	log.Println("Shutting down the server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	web.Shutdown(ctx)

	if os.Getenv("KILL_POSTGRES") == "1" {
		log.Println("Killing postgres...")
		dbutil.KillPostgres()
	}
}
func setupHttpServer(db *pgx.ConnPool, cc *cache.Cache) *http.Server {
	router := httptreemux.New()

	router.POST("/api/forum/create", apiforum.MakeCreateForumHandler(db, cc))
	router.GET("/api/forum/:slug/details", apiforum.MakeShowForumHandler(db))
	router.POST("/api/forum/:slug/create", apiforum.MakeCreateThreadHandler(db, cc))
	router.GET("/api/forum/:slug/users", apiforum.MakeShowUsersHandler(db))
	router.GET("/api/forum/:slug/threads", apiforum.MakeShowThreadsHandler(db))

	router.GET("/api/post/:id/details", apipost.MakeShowPostHandler(db))
	router.POST("/api/post/:id/details", apipost.MakeEditPostHandler(db))

	router.POST("/api/service/clear", apiservice.MakeClearDbHandler(db, cc))
	router.GET("/api/service/status", apiservice.MakeShowStatusHandler(db, cc))

	router.POST("/api/thread/:slug_or_id/create", apithread.MakeCreatePostHandler(db, cc))
	router.GET("/api/thread/:slug_or_id/details", apithread.MakeShowThreadHandler(db))
	router.POST("/api/thread/:slug_or_id/details", apithread.MakeEditThreadHandler(db))
	router.GET("/api/thread/:slug_or_id/posts", apithread.MakeShowPostsHandler(db))
	router.POST("/api/thread/:slug_or_id/vote", apithread.MakeVoteThreadHandler(db))

	router.POST("/api/user/:nickname/create", apiuser.MakeCreateUserHandler(db, cc))
	router.GET("/api/user/:nickname/profile", apiuser.MakeShowUserHandler(db))
	router.POST("/api/user/:nickname/profile", apiuser.MakeEditUserHandler(db))

	if os.Getenv("DEBUG") == "1" {
		fileServ := http.FileServer(http.Dir("."))
		router.GET("/swagger-ui/", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
			fileServ.ServeHTTP(w, r)
		})
		router.GET("/swagger-ui/*path", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
			fileServ.ServeHTTP(w, r)
		})
		router.GET("/", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
			http.Redirect(w, r, "/swagger-ui", http.StatusMovedPermanently)
		})
	}

	handler := negroni.New()
	if os.Getenv("DEBUG") == "1" {
		handler.Use(negroni.NewLogger())
	}
	handler.UseHandler(router)

	server := &http.Server{Addr: ":5000", Handler: handler}
	return server
}
