package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox.xyh.net/internal/models"
	"snippetbox.xyh.net/internal/validator"
	"strconv"
)

// this is used to represent the form data to be sent back to the user in case of a invalid field entry
// use uppercase to let the html template render it
// we also used tags to tell decoder how to map HTML form values into different struct fields

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"` //the decoder will also automatically convert the type to int in this case
	validator.Validator `form:"-"`
}

// the signature of the home handler specifies it is a method of the dependency struct *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//this url checking is not needed anymore since httprouter matches this exactly
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

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
	//we can use ParamsFromContext() and context() to get the parameters needed
	params := httprouter.ParamsFromContext(r.Context())

	//now we get the id value from the parameter
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template. // Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the // snippet expiry to 365 days.
	data.Form = snippetCreateForm{Expires: 365}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//declare a new instance of the instance
	form := snippetCreateForm{}
	//the decodePostForm helper function helps with parse and decode forms
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//we already has a embeded struct Validator in "form" based on its definition
	//we can use the method of Validator directly
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	//if there is invalid entry we need to display it
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
