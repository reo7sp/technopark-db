package main

import (
	//"fmt"
	//"github.com/dimfeld/httptreemux"
	////"github.com/reo7sp/technopark-db/api/forum"
	////"github.com/reo7sp/technopark-db/api/thread"
	//"github.com/reo7sp/technopark-db/dbutil"
	//"log"
	//"net/http"
)

func main() {
	//db, err := dbutil.Connect()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//router := httptreemux.New()
	////router.POST("/api/forum/create", forum.MakeCreateForumHandler(db))
	////router.GET("/api/forum/:slug/details", forum.MakeShowDetailsHandler(db))
	////router.POST("/api/forum/:slug/create", forum.MakeCreateThreadHandler(db))
	////router.POST("/api/thread/:slug/create", thread.MakeCreatePostHandler(db))
	////router.GET("/api/thread/:slug/details", thread.MakeDetailsHandler(db))
	//
	//fileServ := http.FileServer(http.Dir("swagger-ui-dist"));
	//router.GET("/swagger", func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	//	fileServ.ServeHTTP(w, r)
	//})
	//
	//fmt.Println("Starting http server: http://localhost:5000")
	//log.Fatal(http.ListenAndServe(":5000", router))
}
