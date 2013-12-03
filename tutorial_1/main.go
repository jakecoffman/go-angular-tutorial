package main

import (
	"log"
	"net/http"
)

const port = ":80"

func main() {
	// Handle all requests by serving a file of the same name
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("Running on port " + port)
	// This call "blocks" meaning the progam runs here forever
	http.ListenAndServe(port, nil)
}
