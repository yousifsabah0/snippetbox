package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yousifsabah0/snippetbox/pkg/models/mysql"
)

type Application struct {
	errorLogger *log.Logger
	infoLogger *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main () {
	// Get command line flags
	addr := flag.String("addr", ":8080", "HTTP network address.")
	dns := flag.String("dns", "stark:1538@/snippetbox?parseTime=true", "Mysql connection string.")
	flag.Parse()

	// Create custom loggers
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize db connection
	db, err := openDB(*dns)
	if err != nil {
		errorLogger.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLogger.Fatal(err)
	}

	// Initialize new instance of Application struct
	app := &Application{
		errorLogger: errorLogger,
		infoLogger: infoLogger,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize new http.Server struct 
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLogger,
		Handler: app.routes(),
	}

	infoLogger.Printf("Server has been lift off at %v", *addr)
	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
}

func openDB (dns string) (*sql.DB, error) {
	return sql.Open("mysql", dns)
}