package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type handlerError struct {
	Error   error
	Message string
	Code    int
}

type handler func(w http.ResponseWriter, r *http.Request) *handlerError

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("%v\n", err.Error)
		http.Error(w, err.Message, err.Code)
	}
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func listEntries(w http.ResponseWriter, r *http.Request) *handlerError {
	files, e := filepath.Glob("web/blog/*")
	if e != nil {
		return &handlerError{e, "Error getting entries", http.StatusInternalServerError}
	}

	for i, value := range files {
		temp := strings.Replace(value, "\\", "/", -1)
		files[i] = strings.Replace(temp, "web/", "", 1)
	}

	bytes, e := json.Marshal(files)
	if e != nil {
		return &handlerError{e, "Error marshalling JSON", http.StatusInternalServerError}
	}

	w.Write(bytes)
	return nil
}

func main() {
	// command line flags
	port := *flag.Int("port", 80, "port to serve on")
	dir := *flag.String("directory", "web/", "directory of web files")

	// handle all requests by serving a file of the same name
	fs := http.Dir(dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.Handle("/blog", handler(listEntries))

	log.Printf("Running on port %d\n", port)

	// this call blocks -- the progam runs here forever
	portString := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(portString, nil)
	fmt.Println(err.Error())
}
