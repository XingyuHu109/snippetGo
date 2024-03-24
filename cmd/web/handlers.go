package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.xyh.net/internal/models"
	"strconv"
)

// the signature of the home handler specifies it is a method of the dependency struct *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	//retrieve the last 10 snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//use the helper function to create a struct for holing data that include the current year
	data := app.newTemplateData(r)
	data.Snippets = snippets
	//we can create the map of the templates once in main.go using the newTemplateCache() in template.go
	//and then use the render() in helpers.go to execute the chosen template
	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	//use SnippetModel object's Get() to  retrieve the data for a specific record based on its id. If no matching record is found, return 404 response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//use the helper function to create a struct for holing data that include the current year
	data := app.newTemplateData(r)
	data.Snippet = snippet
	//we can create the map of the templates once in main.go using the newTemplateCache() in template.go
	//and then use the render() in helpers.go to execute the chosen template
	app.render(w, http.StatusOK, "view.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	//some dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//redirect the user to  the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
