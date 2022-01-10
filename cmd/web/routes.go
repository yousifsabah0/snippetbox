package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
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
	dynamicMiddlewares := alice.New(app.session.Enable, noSurf)

	mux := pat.New()

	mux.Get("/", dynamicMiddlewares.ThenFunc(app.home))
	mux.Get("/snippets/new", dynamicMiddlewares.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippets/new", dynamicMiddlewares.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippets/:id", dynamicMiddlewares.ThenFunc(app.showSnippet))

	mux.Get("/users/register", dynamicMiddlewares.ThenFunc(app.registerForm))
	mux.Get("/users/login", dynamicMiddlewares.ThenFunc(app.loginForm))

	mux.Post("/users/register", dynamicMiddlewares.ThenFunc(app.register))
	mux.Post("/users/login", dynamicMiddlewares.ThenFunc(app.login))
	mux.Post("/users/logout", dynamicMiddlewares.Append(app.requireAuthentication).ThenFunc(app.logout))

	// Handle static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileserver))

	return middlewares.Then(mux)
}
