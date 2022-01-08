package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes () http.Handler {
	// // Initialize new servemux.
	// mux := http.NewServeMux()

	// // Handle routes
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippets", app.showSnippet)
	// mux.HandleFunc("/snippets/new", app.createSnippet)

	// // Handle static files
	// fileserver := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	middlewares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	// mux.Get("/snippets/new", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippets/new", http.HandlerFunc(app.createSnippet))
	
	mux.Get("/snippets/:id", http.HandlerFunc(app.showSnippet))

		// Handle static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileserver))

	return middlewares.Then(mux)
}