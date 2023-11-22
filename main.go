package main

import (
	"net/http"
	"tp_go/middleware"
	"tp_go/router"
)

func main() {
	const filePath = "dictionary.txt"

	r := router.InitRouter(filePath)

	r.Use(middleware.LoggingMiddleware)

	http.ListenAndServe(":8080", r)
}
