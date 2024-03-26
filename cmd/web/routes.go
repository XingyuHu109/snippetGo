package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	//mux := http.NewServeMux()
	//
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.snippetView)
	//mux.HandleFunc("/snippet/create", app.snippetCreate)
	//
	////use an external package alice to manage the chain for cleaner code
	//standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	////we need to add the CSP header for every request, and this means the middleware function need to wrap around our mux(to execute before it)
	////which middleware comes first is processed first, but also mind the logic flow(in middleware.go)
	//return standard.Then(mux)

	router := httprouter.New()

	//setting a custom handler function in case a file is not found using our router, it acts as a wrapper around the notFound()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { app.notFound(w) })

	//route for the static files
	//the fileserver is a fileHandler in itself
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//then create the routers using the appropriate methods, patterns and handlers
	//the advanced routing already takes care of differentiating between GET and POST requests
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create the middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it
	return standard.Then(router)
}
