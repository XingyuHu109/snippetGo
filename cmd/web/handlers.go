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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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
	//get userId from session data
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	snippets, err := app.snippets.Latest(userID)
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
	//get userID from session data
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	//use SnippetModel object's Get() to  retrieve the data for a specific record based on its id. If no matching record is found, return 404 response
	snippet, err := app.snippets.Get(id, userID)
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

	//flash message is automatically added in the newTemplateDate() function if it exists in the session data

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

	//we already has an embeded struct Validator in "form" based on its definition
	//we can use the method of Validator directly
	//
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
	//get userID from session data
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	//id, err := app.snippets.Insert(form.Title, form.Content, form.Expires, userID)
	_, err = app.snippets.Insert(form.Title, form.Content, form.Expires, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully!")

	//http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)

}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm
	//the definition of userSignupForm Struct allows form.Decoder() to work
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	//if there is any error in inputs, we need to redisplay the page with a 422 code
	//however, unlike the error in creating snippets, we are not re-displaying the password
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	//if no error in input
	//check if the email is duplicate from the db entry
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already in use")
			//error page returned to the user
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
	}

	//otherwise, the signup is successful, and we need to add a flash message to the current session
	app.sessionManager.Put(r.Context(), "flash", "Account successfully created\nPlease log in.")
	//redirect to the login paged
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	//value written to form using decode
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form.
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic // non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//renew the session id(change the session id)
	//it is a good practice to generate a new session once an authentication stage changes(log in or out operation)
	//this retains the data associated with the session, the reason behind this is to prevent session fixation attack
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the ID of the current user to the session, so that they are now // 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//good habit to renew sessions
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	//remove the authenticated user id from the session data to let the user log out
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	//add a flash message to show the user have been logged out
	app.sessionManager.Put(r.Context(), "flash", "Log out successful")

	//redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
