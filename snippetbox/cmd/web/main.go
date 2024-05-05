package main

import (
	"log"
	"net/http"
	// "strconv"
	// "fmt"
)
// func home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Server", "Go")
// 	w.Write([]byte("Hello from Snippetbox"))
// }
// func snippetView(w http.ResponseWriter, r *http.Request) { 
// 	id, err := strconv.Atoi(r.PathValue("id"))
// 	if err != nil || id < 1 {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
//  }

// func snippetCreate(w http.ResponseWriter, r *http.Request) { 
// 	w.Write([]byte("Display a form for creating a new snippet..."))
// }

// func snippetCreatePost(w http.ResponseWriter, r *http.Request) { 
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write([]byte("Save a new snippet..."))
// }
func main() {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/",http.StripPrefix("/static",fileserver))
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView) 
	mux.HandleFunc("GET /snippet/create", snippetCreate) 
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err) 
}