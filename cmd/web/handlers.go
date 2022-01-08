package main

import (
	"errors"
	"fmt"

	// "html/template"
	"net/http"
	"strconv"

	"github.com/yousifsabah0/snippetbox/pkg/models"
)

// Home page handler.
func (app *Application) home (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.FindLatest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// log.Fatal(err.Error())
	// 	// app.errorLogger.Println(err.Error())
	// 	// http.Error(w, "Internal server error.", http.StatusInternalServerError)
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	// log.Fatal(err.Error())
	// 	// app.errorLogger.Println(err.Error())
	// 	// http.Error(w, "Internal server error.", http.StatusInternalServerError)
	// 	app.serverError(w, err)
	// }

	// w.Write([]byte("Hello, World."))
}

// Handler to show specific snippet.
func (app *Application) showSnippet (w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil && id < 1 {
		app.notFound(w)
		return
	}

	// Find snippet
	snippet, err := app.snippets.FindOne(id)
	if err != nil {
		if errors.Is(err, models.ErrNotRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%v", snippet)
}

// Handler to create new snippets.
func (app *Application) createSnippet (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Only 'POST' method allowed.", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Hiiii"
	content := "hii with more words guysss."
	expires := "7"

	id, err := app.snippets.Create(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets?id=%d", id), http.StatusSeeOther)
	// w.Write([]byte("Creating new snippet."))
}