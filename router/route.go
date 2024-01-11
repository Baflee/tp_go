package router

import (
	"fmt"
	"net/http"
	"strconv"
	"tp_go/dictionary"

	"github.com/gorilla/mux"
)

func InitRouter(filePath string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/dictionary/add", func(w http.ResponseWriter, r *http.Request) {
		actionAdd(filePath, w, r)
	}).Methods("POST")

	r.HandleFunc("/dictionary/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(filePath, w, r)
	}).Methods("GET")

	r.HandleFunc("/dictionary/delete/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionRemove(filePath, w, r)
	}).Methods("DELETE")

	r.HandleFunc("/dictionary", func(w http.ResponseWriter, r *http.Request) {
		actionList(filePath, w, r)
	}).Methods("GET")

	return r
}

func actionAdd(filePath string, w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("word")
	value := r.FormValue("definition")

	_, err, statusCode := dictionary.Add(filePath, key, value)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Word %s Added", key)
	}
}

func actionDefine(filePath string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	result, err, statusCode := dictionary.Get(filePath, key)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Definition of %s: %s", key, result)
	}
}

func actionRemove(filePath string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	_, err, statusCode := dictionary.Remove(filePath, key)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Word %s Deleted", key)
	}
}

func actionList(filePath string, w http.ResponseWriter, r *http.Request) {
	result, err, statusCode := dictionary.List(filePath)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Word Lists: \n%s", result)
	}
}
