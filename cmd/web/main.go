package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.xyh.net/internal/models"
	"time"
)

// this is a struct that holds all the application-wide dependencies
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	//creating a command line flag
	// to call the flag use $ go run ./cmd/web/ -addr=":80"
	addr := flag.String("addr", ":4000", "HTTP network address")

	//add a command line flag for the mysql data source name string
	//dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")
	//dsn := flag.String("dsn", "xyh:${DB_PASSWORD}@tcp(snippetapp.mysql.database.azure.com:3306)/snippet?parseTime=true&tls=true", "MySQL data source name")
	////parse the command line flag
	//flag.Parse()

	//create new loggers to separate information and errors.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//putting these 2 log dependencies into the struct

	//create the connection pool
	//db, err := openDB(*dsn)

	password := os.Getenv("DB_PASSWORD")
	expandedDSN := fmt.Sprintf("xyh:%s@tcp(snippetapp.mysql.database.azure.com:3306)/snippet?parseTime=true&tls=true", password)
	db, err := openDB(expandedDSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	//close the connection pool before the main() function is closed
	defer db.Close()

	//initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//initialize a form decoder instance
	formDecoder := form.NewDecoder()

	//initialize a new session manager
	//configure it to use our sql db as the session store, and set a life-time of 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
