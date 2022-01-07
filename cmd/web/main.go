package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main () {
	// Get command line flags
	addr := flag.String("addr", ":8080", "HTTP network address.")
	flag.Parse()

	// Create custom loggers
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize new servemux.
	mux := http.NewServeMux()

	// Handle routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets", showSnippet)
	mux.HandleFunc("/snippets/new", createSnippet)

	// Handle static files
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	infoLogger.Printf("Server has been lift off at %v", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLogger.Fatal(err)
}