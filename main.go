package main

import (
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"github.com/dimfeld/httptreemux"
	"net/http"
	"fmt"
	"github.com/reo7sp/technopark-db/api/apiforum"
	"github.com/reo7sp/technopark-db/api/apithread"
	"github.com/reo7sp/technopark-db/api/apipost"
	"github.com/reo7sp/technopark-db/api/apiservice"
	"github.com/reo7sp/technopark-db/api/apiuser"
)

func main() {
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}

	router := httptreemux.New()

	router.POST("/api/forum/create", apiforum.MakeCreateForumHandler(db))
	router.GET("/api/forum/:slug/details", apiforum.MakeShowForumHandler(db))
	router.POST("/api/forum/:slug/create", apiforum.MakeCreateThreadHandler(db))
	router.POST("/api/forum/:slug/users", apiforum.MakeShowUsersHandler(db))
	router.POST("/api/forum/:slug/threads", apiforum.MakeShowThreadsHandler(db))

	router.GET("/api/post/:id/threads", apipost.MakeShowPostHandler(db))
	router.POST("/api/post/:id/threads", apipost.MakeEditPostHandler(db))

	router.POST("/api/service/clear", apiservice.MakeClearDbHandler(db))
	router.GET("/api/service/status", apiservice.MakeShowStatusHandler(db))

	router.POST("/api/thread/:slug_or_id/create", apithread.MakeCreatePostHandler(db))
	router.GET("/api/thread/:slug_or_id/details", apithread.MakeShowThreadHandler(db))
	router.POST("/api/thread/:slug_or_id/details", apithread.MakeEditThreadHandler(db))
	router.GET("/api/thread/:slug_or_id/posts", apithread.MakeShowPostsHandler(db))
	router.POST("/api/thread/:slug_or_id/vote", apithread.MakeVoteThreadHandler(db))

	router.POST("/api/user/:nickname/create", apiuser.MakeCreateUserHandler(db))
	router.GET("/api/user/:nickname/profile", apiuser.MakeShowUserHandler(db))
	router.POST("/api/user/:nickname/profile", apiuser.MakeEditUserHandler(db))

	fileServ := http.FileServer(http.Dir("."))
	router.GET("/swagger/", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		fileServ.ServeHTTP(w, r)
	})
	router.GET("/swagger/*path", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		fileServ.ServeHTTP(w, r)
	})

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		http.Redirect(w, r, "/swagger", http.StatusMovedPermanently)
	})

	fmt.Println("Starting http server: http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
