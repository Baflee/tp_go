package router

import (
	"fmt"
	"net/http"
	"strconv"
	"tp_go/dictionary"

	"github.com/gorilla/mux"
)

func InitRouter(db string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/dictionary/add", func(w http.ResponseWriter, r *http.Request) {
		actionAdd(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/dictionary/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(db, w, r)
	}).Methods("GET")

	r.HandleFunc("/dictionary/delete/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionRemove(db, w, r)
	}).Methods("DELETE")

	r.HandleFunc("/dictionary", func(w http.ResponseWriter, r *http.Request) {
		actionList(db, w, r)
	}).Methods("GET")

	return r
}

func actionAdd(db string, w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("word")
	value := r.FormValue("definition")

	_, err, statusCode := dictionary.Add(db, "dictionary", key, value)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Word %s Added", key)
	}
}

func actionDefine(db string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	result, err, statusCode := dictionary.Get(db, "dictionary", key)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, result)
	}
}

func actionRemove(db string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	result, err, statusCode := dictionary.Remove(db, "dictionary", key)
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, result)
	}
}

func actionList(db string, w http.ResponseWriter, r *http.Request) {
	result, err, statusCode := dictionary.List(db, "dictionary")
	if err != nil {
		httpErrorMsg := "HTTP Error: " + strconv.Itoa(statusCode) + " - " + err.Error()
		http.Error(w, httpErrorMsg, statusCode)
	} else {
		fmt.Fprintf(w, "Word Lists: \n%s", result)
	}
}
