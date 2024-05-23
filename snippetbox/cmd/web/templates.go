package main

import (
	"snippetbox.jagdish.net/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct{
	Snippet models.Snippet
	Snippets []models.Snippet
	CurrentYear int
	Form any
	Flash string
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04") 
}

var functions = template.FuncMap{
	"humanDate": humanDate, 
}


func newTemplateCache() (map[string]*template.Template, error) { // Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob() function to get a slice of all filepaths that // match the pattern "./ui/html/pages/*.tmpl". This will essentially gives // us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
	return nil, err }
	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
	// Extract the file name (like 'home.tmpl') from the full filepath // and assign it to the name variable.
		name := filepath.Base(page)
		// Create a slice containing the filepaths for our base template, any // partials and the page.
		// files := []string{
		// "./ui/html/base.tmpl.html", "./ui/html/partials/nav.tmpl.html", page,
		// }
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err 
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html") 
		if err != nil {
			return nil, err }
		// Parse the files into a template set.
		ts, err = ts.ParseFiles(page) 
		if err != nil {
			return nil, err 
		}
		// Add the template set to the map, using the name of the page // (like 'home.tmpl') as the key.
		cache[name] = ts
	}
// Return the map.
	return cache, nil 
}