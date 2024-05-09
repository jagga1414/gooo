package main
import ( 
	"fmt"
	"net/http"
	"strconv"
	"errors"
	"snippetbox.jagdish.net/internal/models"

)
func (app *application) home(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Server", "Go")
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

func (app *application)snippetCreate(w http.ResponseWriter, r *http.Request) { 
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) { 
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa" 
	expires := 7
	// Pass the data to the SnippetModel.Insert() method, receiving the // ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}