package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.lilbee/internal/models"
	"snippetbox.lilbee/ui"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
    Form        any
    Flash       string
    IsAuthenticated bool
    CSRFToken   string
    User        models.User
}

 // Create a humanDate function which returns a nicely formatted string
 // representation of a time.Time object.
 func humanDate(t time.Time) string {
   // return t.Format("02 Jan 2006 at 15:04")

    if t.IsZero() {
        return ""
    }
    // Convert the time to UTC before formatting it.
    return t.UTC().Format("02 Jan 2006 at 15:04")
 }
 // Initialize a template.FuncMap object and store it in a global variable. This is
 // essentially a string-keyed map which acts as a lookup between the names of our
 // custom template functions and the functions themselves.
 var functions = template.FuncMap{
    "humanDate": humanDate,
 }

func newTemplateCache() (map[string]*template.Template, error) {
    // Initialize a new map to act as the cache.
    cache := map[string]*template.Template{}
    // Use the filepath.Glob() function to get a slice of all filepaths that
    // match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
    // us a slice of all the filepaths for our application 'page' templates
    // like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
    pages, err := fs.Glob(ui.Files, "html/pages/*.html")
    if err != nil {
        return nil, err
    }
    // Loop through the page filepaths one-by-one.
    for _, page := range pages {
        // Extract the file name (like 'home.tmpl') from the full filepath
        // and assign it to the name variable.
        name := filepath.Base(page)

        patterns := []string{
            "html/base.html",
            "html/partials/*.html",
            page,
        }

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
    
        // Add the template set to the map, using the name of the page
        // (like 'home.tmpl') as the key.
        cache[name] = ts
    }
    // Return the map.
    return cache, nil
 }