package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

func ReadJsonObject(w http.ResponseWriter, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	return nil
}

func WriteJsonObject(w http.ResponseWriter, v interface{}, statusCode int) error {
	respBody, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(500)
		log.Println("error: api.WriteJsonObject: json.Marshal:", err)
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)
	return nil
}
