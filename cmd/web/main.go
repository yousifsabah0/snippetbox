package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLogger *log.Logger
	infoLogger *log.Logger
}

func main () {
	// Get command line flags
	addr := flag.String("addr", ":8080", "HTTP network address.")
	flag.Parse()

	// Create custom loggers
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize new instance of Application struct
	app := &Application{
		errorLogger: errorLogger,
		infoLogger: infoLogger,
	}


	// Initialize new http.Server struct 
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLogger,
		Handler: app.routes(),
	}

	infoLogger.Printf("Server has been lift off at %v", *addr)
	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
}