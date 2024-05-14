package main

import ( 
	"log"
	"net/http"
)
func home(w http.ResponseWriter, r *http.Request) { 
	print("why")
	w.Write([]byte("Hello from Snippetbox"))
}
// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) { 
	w.Write([]byte("Display a specific snippet..."))
}
// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Display a form for creating a new snippet..."))
}
func main() {
// Register the two new handler functions and corresponding route patterns with // the servemux, in exactly the same way that we did before.
mux := http.NewServeMux()

mux.HandleFunc("GET /comment", home)
mux.HandleFunc("/snippet/view/{id}", snippetView)
mux.HandleFunc("/snippet/create", snippetCreate)

log.Print("starting server on :8080")
err := http.ListenAndServe(":8080", mux)
log.Fatal(err) 
}