package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const blog = "blog/"

var router = mux.NewRouter()

type handlerError struct {
	Error   error
	Message string
	Code    int
}

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

func listEntries(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	files, e := filepath.Glob(blog + "*")
	if e != nil {
		return nil, &handlerError{e, "Error getting entries", http.StatusInternalServerError}
	}

	entries := make([]map[string]string, 0, len(files))
	for _, value := range files {
		name := strings.Replace(value, "\\", "/", -1) // in case of windows
		name = strings.Replace(name, blog, "", 1)
		entry := map[string]string{
			"name": name,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func getEntry(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	entry := mux.Vars(r)["entry"]
	content, e := ioutil.ReadFile(blog + entry)
	if e != nil {
		return nil, &handlerError{e, "Could not find entry", http.StatusNotFound}
	}

	response := map[string]string{
		"content": string(content),
	}

	return response, nil
}

func addEntry(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var payload map[string]string
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}
	if !containsAll(payload, "name", "content") {
		e = errors.New("JSON must contain 'name' and 'content'")
		return nil, &handlerError{e, e.Error(), http.StatusBadRequest}
	}

	e = ioutil.WriteFile(blog+payload["name"], []byte(payload["content"]), 0644)
	if e != nil {
		return nil, &handlerError{e, "Could not write to disk", http.StatusInternalServerError}
	}

	return make(map[string]string), nil
}

func containsAll(m map[string]string, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	router.Handle("/", fileHandler)
	router.Handle("/static", fileHandler)
	router.Handle("/blog", handler(listEntries))
	router.Handle("/blog/{entry}", handler(addEntry)).Methods("POST")
	router.Handle("/blog/{entry}", handler(getEntry)).Methods("GET")
	http.Handle("/", router)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
