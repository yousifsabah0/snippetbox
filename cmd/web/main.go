package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/yousifsabah0/snippetbox/pkg/models/mysql"
)

type Application struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
}

func main() {
	// Get command line flags
	addr := flag.String("addr", ":8080", "HTTP network address.")
	dns := flag.String("dns", "stark:1538@/snippetbox?parseTime=true", "Mysql connection string.")
	secret := flag.String("secret", "eef57723d3ff31cc2a0033289172ff39", "Secret session key")

	// Parse flags
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

	// Initialize session
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Initialize new instance of Application struct
	app := &Application{
		errorLogger:   errorLogger,
		infoLogger:    infoLogger,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize new http.Server struct
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLogger,
		Handler:      app.routes(),
		TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Server has been lift off at %v", *addr)
	err = srv.ListenAndServeTLS("./certificate/cert.pem", "./certificate/key.pem")
	errorLogger.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	return sql.Open("mysql", dns)
}
