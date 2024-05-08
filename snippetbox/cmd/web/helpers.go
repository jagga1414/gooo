package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri,"trace",trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) 
}
	// The clientError helper sends a specific status code and corresponding description // to the user. We'll use this later in the book to send responses like 400 "Bad
	// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status) 
}


func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) { // Retrieve the appropriate template set from the cache based on the page
		// name (like 'home.tmpl'). If no entry exists in the cache with the
		// provided name, then create a new error and call the serverError() helper
		// method that we made earlier and return.
		buf := new(bytes.Buffer)
		ts, ok := app.templateCache[page] 
		if !ok {
			err := fmt.Errorf("the template %s does not exist", page) 
			app.serverError(w, r, err)
			return
		}
		// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
		
		// Execute the template set and write the response body. Again, if there // is any error we call the serverError() helper.
		err := ts.ExecuteTemplate(buf, "base", data)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		w.WriteHeader(status)

		buf.WriteTo(w)
	}