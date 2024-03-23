package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// this is a struct that holds all the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//creating a command line flag
	// to call the flag use $ go run ./cmd/web/ -addr=":80"
	addr := flag.String("addr", ":4000", "HTTP network address")
	//parse the command line flag
	flag.Parse()

	//create new loggers to separate information and errors.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//putting these 2 log dependencies into the struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
