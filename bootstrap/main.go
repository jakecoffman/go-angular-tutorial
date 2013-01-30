package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"os"
)

const port = ":8090"

var myTemplate = template.Must(template.ParseFiles(
	"templates/index.html",
))

type MyResponse struct{
	Title string
	Page string
}

func myHandler(w http.ResponseWriter, r *http.Request){
	response := MyResponse{}
	response.Title = "My title"
	response.Page = mux.Vars(r)["page"]
	if err := myTemplate.Execute(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	r := mux.NewRouter()

	r.HandleFunc("/", myHandler)
	r.HandleFunc("/{page}", myHandler)
	
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	fmt.Println("Running on " + port)
	http.ListenAndServe(port, nil)
}
