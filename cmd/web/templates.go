package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.xyh.net/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		//Parse the base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.html")
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
