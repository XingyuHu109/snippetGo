package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.xyh.net/internal/models"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
}

// returns a nicely formated time
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// create a template.FuncMap object, this is basically a lookup map that helps us locate the right function name
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	//new map to act as the cache
	cache := map[string]*template.Template{}

	//get all the file path that matches the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	//loop through thr pages file path 1 through 1
	for _, page := range pages {
		name := filepath.Base(page)

		// this is the way of registering a function in a template
		// The template.FuncMap must be registered with the template set before you // call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		//call ParseGlob() on this template set(ts) to add all the partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		//next, we need to parse the current page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//add the template set to the map, using the base name of the page
		cache[name] = ts
	}
	return cache, nil
}
