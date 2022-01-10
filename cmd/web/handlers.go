package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yousifsabah0/snippetbox/pkg/forms"
	"github.com/yousifsabah0/snippetbox/pkg/models"
)

// =========================================
// Snippets handlers
// =========================================

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

	app.render(w, r, "home.page.tmpl", &TemplateData{
		Snippets: snippets,
	})
}

// Handler to show specific snippet.
func (app *Application) showSnippet (w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	app.render(w, r, "show.page.tmpl", &TemplateData{
		Snippet: snippet,
	})
}

// Handler to show new snippet form
func (app *Application) createSnippetForm (w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

// Handler to create new snippets.
func (app *Application) createSnippet (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 50)
	form.PremittedValue("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &TemplateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Create(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "New snippet created.")
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

// =========================================
// Users handlers
// =========================================

// Display register page
func (app *Application) registerForm (w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

// Display login page
func (app *Application) loginForm (w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

// Register endpoint
func (app *Application) register (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)

	form.Required("name", "email", "password")
	
	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &TemplateData{
			Form: form,
		})
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address already exists")
			app.render(w, r, "register.page.tmpl", &TemplateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Account created. Go and login")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

// Login endpoint
func (app *Application) login (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)

	form.Required("email", "password")
	if !form.Valid() {
		app.render(w, r, "login.page.tmpl", &TemplateData{
			Form: form,
		})
	}

	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password is incorrect")
			app.render(w, r, "login.page.tmpl", &TemplateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserId", id)
	http.Redirect(w, r, "/snippets/new", http.StatusSeeOther)
}

func (app *Application) logout (w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserId")
	app.session.Put(r, "flash", "You have been loggedout")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}