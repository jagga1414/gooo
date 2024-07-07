package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description // to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}


func (app *application) renderV(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) { // Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	// buf := new(bytes.Buffer)
	// ts, ok := app.templateCache[page]
	// if !ok {
	// 	err := fmt.Errorf("the template %s does not exist", page)
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// // Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).

	// var err error
	// ts1, err := ts.ParseGlob("./ui/html/blog_pages/me.tmpl.html") 
	// if err != nil {
	// 	// err = fmt.Errorf("template not found")
	// 	app.serverError(w,r,err)
	// 	return }

	// // Execute the template set and write the response body. Again, if there // is any error we call the serverError() helper.
	// err = ts1.ExecuteTemplate(buf, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// w.WriteHeader(status)

	// buf.WriteTo(w)


	files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/pages/view.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/partials/footer.tmpl.html",
			}
			// Use the template.ParseFiles() function to read the files and store the // templates in a template set. Notice that we use ... to pass the contents // of the files slice as variadic arguments.
	ts, err := template.New("me").Funcs(functions).ParseFiles(files...)
	if err != nil {
		app.serverError(w,r,err) 
		return
	}

	ts,err = ts.ParseGlob("./ui/html/blog_pages/"+data.Snippet.FileName)
	// Use the ExecuteTemplate() method to write the content of the "base" // template as the response body.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w,r,err) }
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

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) isAdmin(r *http.Request) bool {
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok {
		return false
	}
	return isAdmin
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		IsAdmin: app.isAdmin(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our // snippetCreatePost handler.
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
