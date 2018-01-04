package service

import (
	"database/sql"
	"net/http"
	"log"
	"github.com/reo7sp/technopark-db/apiutil"
)

func MakeClearDbHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		err := ClearDB(db)
		if err != nil {
			log.Println("error: api.service.MakeClearDbHandler: ClearDB:", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	}
}

func MakeShowStatusHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		statusModel, err := LoadStatusFromDB(db)
		if err != nil {
			log.Println("error: api.service.MakeShowStatusHandler: LoadStatusFromDB:", err)
			w.WriteHeader(500)
			return
		}

		apiutil.WriteJsonObject(w, statusModel, 200)
	}
}
