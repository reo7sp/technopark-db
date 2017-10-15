package forum

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"github.com/reo7sp/technopark-db/database"
)

func CreateFuncMaker(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panic(err)
		}

		var forum Forum
		err = json.Unmarshal(body, &forum)
		if err != nil {
			log.Panic(err)
		}

		err = forum.Create(db)
		isDublicate := false
		if err != nil {
			if database.IsErrorAboutDublicate(err) {
				isDublicate = true
				err = forum.FetchPostsAndThreadsCount(db)
				if err != nil {
					log.Panic(err)
				}
			} else {
				log.Panic(err)
			}
		}

		respBody, err := json.Marshal(forum)
		if err != nil {
			log.Panic(err)
		}
		if isDublicate {
			w.WriteHeader(409)
		} else {
			w.WriteHeader(201)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)
	}
}
