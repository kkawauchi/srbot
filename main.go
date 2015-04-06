// +build !appengine

package main

import (
	"log"
	"net/http"

	_ "./bot"
)

func main() {
	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
