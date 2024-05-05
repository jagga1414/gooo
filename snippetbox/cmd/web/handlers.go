package main
import ( 
	"fmt"
	"net/http"
	"strconv"
	"log"
	"html/template"

)
func home(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Server", "Go")
		// Use the template.ParseFiles() function to read the template file into a
		// template set. If there's an error, we log the detailed error message, use
		// the http.Error() function to send an Internal Server Error response to the 
		// user, and then return from the handler so no subsequent code is executed. 
		files := []string{"./ui/html/pages/home.tmpl.html","./ui/html/pages/base.tmpl.html","./ui/html/partials/nav.tmpl.html"}
		ts, err := template.ParseFiles(files...)
		if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		 return
		}
		// Then we use the Execute() method on the template set to write the
		// template content as the response body. The last parameter to Execute() // represents any dynamic data that we want to pass in, which for now we'll // leave as nil.
		err = ts.ExecuteTemplate(w,"base",nil)
		if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError) 
	}
}
func snippetView(w http.ResponseWriter, r *http.Request) { 
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
 }

func snippetCreate(w http.ResponseWriter, r *http.Request) { 
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) { 
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}