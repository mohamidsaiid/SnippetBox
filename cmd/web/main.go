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

	"mohamidsaiid.com/snippetbox/pkg/models"
	"mohamidsaiid.com/snippetbox/pkg/models/mysql"

	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

// the main dependency used through the application that contains some errors and its own methods
// as handlers and helper and routes
type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	snippets      interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {

	// addr should contain which port we would be working on
	// the user could pass is as input before running the server
	// if not it would contain the port :4000
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	// data base user and his own password
	// could be passed as input through the cli or it would have the it own fixed value
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	// for getting  the user input through the cli
	flag.Parse()

	// creats my own form of error by using the log.New prop
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// intializing your template cache object
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))

	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// initalizing your own app object
	app := &Application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// overloading the http.Server struct from the net/http package to take my own
	// 1- add   	2- handler	3- errorlog type
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Println("Starting server on port ", *addr)

	// this line starts your own server to get and send response
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	// to start your own database pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// ping to check the connection if estaplished or not
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// return your own database connection
	return db, nil
}
