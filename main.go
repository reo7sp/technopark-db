package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"github.com/reo7sp/technopark-db/api/forum"
	"github.com/reo7sp/technopark-db/database"
	"fmt"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.POST("/api/forum/create", forum.CreateFuncMaker(db))
	router.GET("/api/forum/:slug/details", forum.DetailsFuncMaker(db))
	router.POST("/api/forum/:slug/create", forum.CreateThreadFuncMaker(db))
	router.NotFound = http.FileServer(http.Dir("swagger-ui-dist"))

	fmt.Println("Starting http server at 5000 port")
	log.Fatal(http.ListenAndServe(":5000", router))
}
