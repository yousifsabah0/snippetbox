package main

import (
	"net/http"
)

func (app *Application) routes () *http.ServeMux {
	// Initialize new servemux.
	mux := http.NewServeMux()

	// Handle routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippets", app.showSnippet)
	mux.HandleFunc("/snippets/new", app.createSnippet)

	// Handle static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	return mux
}