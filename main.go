package main

import (
	"fmt"
	"net/http"
	"tp_go/dictionary"

	"github.com/gorilla/mux"
)

func main() {
	const filePath = "dictionary.txt"
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

	http.ListenAndServe(":8080", r)
}

func actionAdd(filePath string, w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("word")
	value := r.FormValue("definition")

	_, err := dictionary.Add(filePath, key, value)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "Word %s Added", key)
	}
}

func actionDefine(filePath string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	result, err := dictionary.Get(filePath, key)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusNotFound)
	} else {
		fmt.Fprintf(w, "Definition of %s: %s", key, result)
	}
}

func actionRemove(filePath string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["word"]

	_, err := dictionary.Remove(filePath, key)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "Word %s Deleted", key)
	}
}

func actionList(filePath string, w http.ResponseWriter, r *http.Request) {
	result, err := dictionary.List(filePath)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "Word Lists: \n%s", result)
	}
}
