package main

import (
	"log"
	"net/http"
)

/*
	Calls the http GET/POST handler function at the specified URL
*/
func main() {
	http.HandleFunc("/messages/messages", handleMessages)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
