package apiutil

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ReadJsonObject(r *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}

	return nil
}

func WriteJsonObject(w http.ResponseWriter, obj interface{}, statusCode int) error {
	respBody, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(500)
		log.Println("error: apiutil.WriteJsonObject: json.Marshal:", err)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(respBody)
	return nil
}
