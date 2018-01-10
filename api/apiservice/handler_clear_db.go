package apiservice

import (
	"log"
	"net/http"
	"github.com/jackc/pgx"
)

func MakeClearDbHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		clearDbAction(w, db)
	}
	return f
}

func clearDbAction(w http.ResponseWriter, db *pgx.ConnPool) {
	_, err := db.Exec("TRUNCATE TABLE forums, threads, users, posts RESTART IDENTITY CASCADE")
	if err != nil {
		log.Println("error: apiservice.clearDbAction: TRUNCATE:", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
