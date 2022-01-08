package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *Application) addDefaultDate (td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	td.CurrentYear = time.Now().Year()
	return td
}

func (app *Application) render (w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s don't exists", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultDate(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}

func (app *Application) serverError (w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s/n%s", err.Error(), debug.Stack())
	app.errorLogger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError (w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound (w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}