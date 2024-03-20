package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //to only include the home page, and show 404 for other wildcard
		http.NotFound(w, r)
	} else {
		w.Write([]byte("hello from SnippetBox")) //this is the respones body
	}

}

func snipperView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Displaying a specific snippet")) //this is the respones body
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet...")) //this is the respones body
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snipperView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("starting server on port 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
