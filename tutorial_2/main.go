package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := *flag.Int("port", 80, "port to serve on")
	dir := *flag.String("directory", "web/", "directory of web files")
	// Handle all requests by serving a file of the same name
	http.Handle("/", http.FileServer(http.Dir(dir)))
	log.Printf("Running on port %d\n", port)
	// This call "blocks" meaning the progam runs here forever
	p := fmt.Sprintf(":%d", port)
	http.ListenAndServe(p, nil)
}
