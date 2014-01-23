package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	handler := http.FileServer(fs)
	http.Handle("/", handler)

	log.Printf("Running on port %d\n", *port)

	host := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(host, nil)
	fmt.Println(err.Error())
}
