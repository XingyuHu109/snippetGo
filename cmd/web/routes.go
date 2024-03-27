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

	//create a new middleware chain containing the middleware specific to our dynamic router(not including the file server, since it does
	//not need to be stateful)
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	//then create the routers using the appropriate methods, patterns and handlers
	//the advanced routing already takes care of differentiating between GET and POST requests
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Add the five new routes, all of which use our 'dynamic' middleware chain.
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// Create the middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Wrap the router with the middleware and return it
	return standard.Then(router)
}
