package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Home page handler.
func home (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello, World."))
}

// Handler to show specific snippet.
func showSnippet (w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Showing snippet with id %v", id)
}

// Handler to create new snippets.
func createSnippet (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Only 'POST' method allowed.", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Creating new snippet."))
}