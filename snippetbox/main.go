package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("hello I am root"))
}

func snippetView(w http.ResponseWriter, r *http.Request){
	viewId,err := strconv.Atoi(r.PathValue("id"))	

	if (err!=nil || viewId<1){
		http.NotFound(w,r)
	}
	msg := fmt.Sprintf("I am a snippet %d ...",viewId)
	w.Write([]byte(msg))
}
func snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("create snippet"))
}

func main(){

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}",home)
	mux.HandleFunc("/snippet/view/{id}",snippetView)

	mux.HandleFunc("/snippet/create",snippetCreate)
	log.Print("server started at 4000")

	err := http.ListenAndServe(":4000",mux)

	log.Fatal(err)
}
