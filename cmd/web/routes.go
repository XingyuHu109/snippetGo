package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	//use an external package alice to manage the chain for cleaner code
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//we need to add the CSP header for every request, and this means the middleware function need to wrap around our mux(to execute before it)
	//which middleware comes first is processed first, but also mind the logic flow(in middleware.go)
	return standard.Then(mux)
}
