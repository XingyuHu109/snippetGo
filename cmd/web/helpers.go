package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// retrieve the appropriate template from the cache based on the name of the page
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	//initialize a new buffer used for pre-rendering the html page
	buf := new(bytes.Buffer)

	//write the template to the buffer, instead f straight to the http.response writer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//write the provided http status code
	w.WriteHeader(status)

	//execute the template retrieved from the map, with the data being put into the template
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Add the flash message to the template data, if one exists.
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst, // is the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our // createSnippetPost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as // the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method // will return an error with the type *form.InvalidDecoderError.We use // errors.As() to check for this and raise a panic rather than returning // the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err
	}
	return nil
}

// check if the request is from an authenticated user by checking their session data
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
