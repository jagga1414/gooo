package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.jagdish.net/internal/models"
	"snippetbox.jagdish.net/internal/validator"
)


type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
	}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

		// w.Header().Add("Server", "Go")
		// Use the template.ParseFiles() function to read the template file into a
		// template set. If there's an error, we log the detailed error message, use
		// the http.Error() function to send an Internal Server Error response to the 
		// user, and then return from the handler so no subsequent code is executed. 
	snippets, err := app.snippets.Latest() 
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets 
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}
func (app *application)snippetView(w http.ResponseWriter, r *http.Request) { 
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) { 
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
// Write the snippet data as a plain-text HTTP response body.
	// fmt.Fprintf(w, "%+v", snippet)
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)

 }
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) { 
	data := app.newTemplateData(r)
// Initialize a new createSnippetForm instance and pass it to the template. // Notice how this is also a great opportunity to set any default or
// 'initial' values for the form --- here we set the initial value for the // snippet expiry to 365 days.
	data.Form = snippetCreateForm{ 
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}



func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) { 
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm
	// app.formDecoder.Decode(&form,r.PostForm)
	err = app.decodePostForm(r,&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Because the Validator struct is embedded by the snippetCreateForm struct, 
	// we can call CheckField() directly on it to execute our validation checks. 
	// CheckField() will add the provided key and error message to the
	// FieldErrors map if the check does not evaluate to true. For example, in 
	// the first line here we "check that the form.Title field is not blank". In 
	// the second, we "check that the form.Title field has a maximum character 
	// length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank") 
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long") 
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank") 
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	// Use the Valid() method to see if any of the checks failed. If they did, // then re-render the template passing in the form in the same way as
	// before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data) 
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires) 
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther) 
}