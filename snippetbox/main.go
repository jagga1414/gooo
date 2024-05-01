package main


import (
    "log"
    "net/http"
)

func home(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("hello I am root"))
}

func snippetView(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("I am a snippet"))
}
func snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("create snippet"))
}

func main(){

	mux := http.NewServeMux()

	mux.HandleFunc("/",home)
	mux.HandleFunc("/snippet/view",snippetView)

	mux.HandleFunc("/snippet/create",snippetCreate)
	log.Print("server started at 4000")

	err := http.ListenAndServe(":4000",mux)

	log.Fatal(err)
}
