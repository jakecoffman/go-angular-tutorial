package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const blog = "blog/"

var router = mux.NewRouter()

type handlerError struct {
	Error   error
	Message string
	Code    int
}

type book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Id     int    `json:"id"`
}

var books = make([]book, 0)

type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := fn(w, r)
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}

	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func listBooks(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	return books, nil
}

func getBook(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	b, index := getBookById(id)

	if index < 0 {
		return nil, &handlerError{nil, "Could not find book " + param, http.StatusNotFound}
	}

	return b, nil
}

func addBook(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var payload book
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	payload.Id = getNextId()
	books = append(books, payload)

	return make(map[string]string), nil
}

func removeBook(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	_, index := getBookById(id)

	if index < 0 {
		return nil, &handlerError{nil, "Could not find entry " + param, http.StatusNotFound}
	}

	books = books[:index+copy(books[index:], books[index+1:])]
	return make(map[string]string), nil
}

func getBookById(id int) (book, int) {
	for i, b := range books {
		if b.Id == id {
			return b, i
		}
	}
	return book{}, -1
}

var id = 0

func getNextId() int {
	id += 1
	return id
}

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	router.Handle("/", http.RedirectHandler("/static/", 302))
	router.Handle("/books", handler(listBooks)).Methods("GET")
	router.Handle("/books", handler(addBook)).Methods("POST")
	router.Handle("/books/{id}", handler(getBook)).Methods("GET")
	router.Handle("/books/{id}", handler(removeBook)).Methods("DELETE")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
	http.Handle("/", router)

	books = append(books, book{"Ender's Game", "Orson Scott Card", getNextId()})
	books = append(books, book{"Code Complete", "Steve McConnell", getNextId()})
	books = append(books, book{"World War Z", "Max Brooks", getNextId()})
	books = append(books, book{"Pragmatic Programmer", "David Thomas", getNextId()})

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
