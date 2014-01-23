package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// command line flags
	port := *flag.Int("port", 80, "port to serve on")
	dir := *flag.String("directory", "web/", "directory of web files")

	// handle all requests by serving a file of the same name
	fs := http.Dir(dir)
	handler := http.FileServer(fs)
	http.Handle("/", handler)

	log.Printf("Running on port %d\n", port)

	portString := fmt.Sprintf(":%d", port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(portString, nil)
	fmt.Println(err.Error())
}
