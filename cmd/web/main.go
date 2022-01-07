package main

import (
	"log"
	"net/http"
)

func main () {
	mux := http.NewServeMux()

	// Handle routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets", showSnippet)
	mux.HandleFunc("/snippets/new", createSnippet)

	// Handle static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", fileserver)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}